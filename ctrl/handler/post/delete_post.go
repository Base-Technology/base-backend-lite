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

func DeletePostHandle(c *gin.Context) {
	hd := &DeletePostHandler{}
	handler.Handle(c, hd)
}

type DeletePostHandler struct {
	Req  DeletePostRequest
	Resp DeletePostResponse
}

type DeletePostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	User   *database.User
}

type DeletePostResponse struct {
	common.BaseResponse
}

func (h *DeletePostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *DeletePostHandler) AfterBindReq() error {
	return nil
}

func (h *DeletePostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *DeletePostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *DeletePostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *DeletePostHandler) NeedVerifyToken() bool {
	return true
}

func (h *DeletePostHandler) Process() {
	if err := database.GetInstance().Where("id = ?", h.Req.PostID).Where("creator_id = ?", h.Req.User.ID).Unscoped().Delete(&database.Post{}).Error; err != nil {
		msg := fmt.Sprintf("delete post error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
