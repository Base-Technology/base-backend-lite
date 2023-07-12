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

func CollectPostHandle(c *gin.Context) {
	hd := &CollectPostHandler{}
	handler.Handle(c, hd)
}

type CollectPostHandler struct {
	Req  CollectPostRequest
	Resp CollectPostResponse
}

type CollectPostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	User   *database.User
}

type CollectPostResponse struct {
	common.BaseResponse
}

func (h *CollectPostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *CollectPostHandler) AfterBindReq() error {
	return nil
}

func (h *CollectPostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *CollectPostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *CollectPostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *CollectPostHandler) NeedVerifyToken() bool {
	return true
}

func (h *CollectPostHandler) Process() {
	collect := &database.Collect{UserID: h.Req.User.ID, PostID: h.Req.PostID}
	if err := database.GetInstance().Create(collect).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
