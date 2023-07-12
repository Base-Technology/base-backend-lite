package post

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func UncollectPostHandle(c *gin.Context) {
	hd := &UncollectPostHandler{}
	handler.Handle(c, hd)
}

type UncollectPostHandler struct {
	Req  UncollectPostRequest
	Resp UncollectPostResponse
}

type UncollectPostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	User   *database.User
}

type UncollectPostResponse struct {
	common.BaseResponse
}

func (h *UncollectPostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *UncollectPostHandler) AfterBindReq() error {
	return nil
}

func (h *UncollectPostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *UncollectPostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *UncollectPostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *UncollectPostHandler) NeedVerifyToken() bool {
	return true
}

func (h *UncollectPostHandler) Process() {
	if err := database.GetInstance().Where("post_id = ?", h.Req.PostID).Where("user_id = ?", h.Req.User.ID).Unscoped().Delete(&database.Collect{}).Error; err != nil {
		msg := fmt.Sprintf("uncollect post error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
