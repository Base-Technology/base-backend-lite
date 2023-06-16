package group

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetGroupMemberHandle(c *gin.Context) {
	hd := &GetGroupMemberHandler{}
	handler.Handle(c, hd)
}

type GetGroupMemberHandler struct {
	Req  GetGroupMemberRequest
	Resp GetGroupMemberResponse
}

type GetGroupMemberRequest struct {
	GroupID uint `json:"groupid"`
	User    *database.User
}

type GetGroupMemberResponse struct {
	common.BaseResponse
	GroupMembers []*MemberDetail `json:"data"`
}
type MemberDetail struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (h *GetGroupMemberHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *GetGroupMemberHandler) AfterBindReq() error {
	return nil
}

func (h *GetGroupMemberHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetGroupMemberHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetGroupMemberHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetGroupMemberHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetGroupMemberHandler) Process() {
	var err error
	var group *database.Group
	if err = database.GetInstance().Model(&database.Group{}).Preload("Members").Where("id = ?", h.Req.GroupID).Find(&group).Error; err != nil {
		msg := fmt.Sprintf("get friend list error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	for _, member := range group.Members {
		h.Resp.GroupMembers = append(h.Resp.GroupMembers, &MemberDetail{
			ID:     member.ID,
			Name:   member.Name,
			Avatar: member.Avatar,
		})
	}
}
