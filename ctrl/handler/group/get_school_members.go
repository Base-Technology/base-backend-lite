package group

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/ctrl/types"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetGroupUserHandle(c *gin.Context) {
	hd := &GetGroupUserHandler{}
	handler.Handle(c, hd)
}

type GetGroupUserHandler struct {
	Req  GetGroupUserRequest
	Resp GetGroupUserResponse
}

type GetGroupUserRequest struct {
	Type  string
	Page  int `json:"page"`
	Limit int `json:"limit"`
	User  *database.User
}

type GetGroupUserResponse struct {
	common.BaseResponse
	Users []*types.UserDetail `json:"data"`
}

func (h *GetGroupUserHandler) BindReq(c *gin.Context) error {
	h.Req.Type = c.Query("type")
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *GetGroupUserHandler) AfterBindReq() error {
	return nil
}

func (h *GetGroupUserHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetGroupUserHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetGroupUserHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetGroupUserHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetGroupUserHandler) Process() {
	users := []*database.User{}
	if err := database.GetInstance().Model(&database.User{}).Offset((h.Req.Page-1)*h.Req.Limit).Limit(h.Req.Limit).Where("school = ?", h.Req.User.School).Find(&users).Error; err != nil {
		msg := fmt.Sprintf("get group user error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Users = []*types.UserDetail{}
	for _, user := range users {
		h.Resp.Users = append(h.Resp.Users, &types.UserDetail{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		})
	}
}
