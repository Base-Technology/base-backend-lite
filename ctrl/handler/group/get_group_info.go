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

func GetGroupInfoHandle(c *gin.Context) {
	hd := &GetGroupInfoHandler{}
	handler.Handle(c, hd)
}

type GetGroupInfoHandler struct {
	Req  GetGroupInfoRequest
	Resp GetGroupInfoResponse
}

type GetGroupInfoRequest struct {
	GroupID uint `json:"groupid"`
	User    *database.User
}

type GetGroupInfoResponse struct {
	common.BaseResponse
	GroupInfo *database.Group `json:"data"`
}

func (h *GetGroupInfoHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *GetGroupInfoHandler) AfterBindReq() error {
	return nil
}

func (h *GetGroupInfoHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetGroupInfoHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetGroupInfoHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetGroupInfoHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetGroupInfoHandler) Process() {
	var err error
	var group *database.Group
	if err = database.GetInstance().Model(&database.Group{}).Preload("Members").Where("id = ?", h.Req.GroupID).Find(&group).Error; err != nil {
		msg := fmt.Sprintf("get friend list error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.GroupInfo = group
}
