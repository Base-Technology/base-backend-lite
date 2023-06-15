package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CancelFollowHandle(c *gin.Context) {
	hd := &CancelFollowHandler{}
	handler.Handle(c, hd)
}

type CancelFollowHandler struct {
	Req  CancelFollowRequest
	Resp CancelFollowResponse
}

type CancelFollowRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	User   *database.User
}

type CancelFollowResponse struct {
	common.BaseResponse
}

func (h *CancelFollowHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *CancelFollowHandler) AfterBindReq() error {
	return nil
}

func (h *CancelFollowHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *CancelFollowHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *CancelFollowHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *CancelFollowHandler) NeedVerifyToken() bool {
	return true
}

func (h *CancelFollowHandler) Process() {
	if err := database.GetInstance().Where("user_id = ?", h.Req.User.ID).Where("following_id = ?", h.Req.UserID).Unscoped().Delete(&database.Follow{}).Error; err != nil {
		msg := fmt.Sprintf("delete follow relation error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
