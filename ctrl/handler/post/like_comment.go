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

func LikeCommentHandle(c *gin.Context) {
	hd := &LikeCommentHandler{}
	handler.Handle(c, hd)
}

type LikeCommentHandler struct {
	Req  LikeCommentRequest
	Resp LikeCommentResponse
}

type LikeCommentRequest struct {
	CommentID uint `json:"comment_id" binding:"required"`
	User      *database.User
}

type LikeCommentResponse struct {
	common.BaseResponse
}

func (h *LikeCommentHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *LikeCommentHandler) AfterBindReq() error {
	return nil
}

func (h *LikeCommentHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *LikeCommentHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *LikeCommentHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *LikeCommentHandler) NeedVerifyToken() bool {
	return true
}

func (h *LikeCommentHandler) Process() {
	likecomment := &database.Likecomment{UserID: h.Req.User.ID, CommentID: h.Req.CommentID}
	if err := database.GetInstance().Create(likecomment).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
