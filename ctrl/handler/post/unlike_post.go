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

func UnlikePostHandle(c *gin.Context) {
	hd := &UnlikePostHandler{}
	handler.Handle(c, hd)
}

type UnlikePostHandler struct {
	Req  UnlikePostRequest
	Resp UnlikePostResponse
}

type UnlikePostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	User   *database.User
}

type UnlikePostResponse struct {
	common.BaseResponse
}

func (h *UnlikePostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *UnlikePostHandler) AfterBindReq() error {
	return nil
}

func (h *UnlikePostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *UnlikePostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *UnlikePostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *UnlikePostHandler) NeedVerifyToken() bool {
	return true
}

func (h *UnlikePostHandler) Process() {
	if err := database.GetInstance().Where("post_id = ?", h.Req.PostID).Where("user_id = ?", h.Req.User.ID).Unscoped().Delete(&database.Like{}).Error; err != nil {
		msg := fmt.Sprintf("unlike post error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
