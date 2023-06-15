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

func LikePostHandle(c *gin.Context) {
	hd := &LikePostHandler{}
	handler.Handle(c, hd)
}

type LikePostHandler struct {
	Req  LikePostRequest
	Resp LikePostResponse
}

type LikePostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
	User   *database.User
}

type LikePostResponse struct {
	common.BaseResponse
}

func (h *LikePostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *LikePostHandler) AfterBindReq() error {
	return nil
}

func (h *LikePostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *LikePostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *LikePostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *LikePostHandler) NeedVerifyToken() bool {
	return true
}

func (h *LikePostHandler) Process() {
	like := &database.Like{UserID: h.Req.User.ID, PostID: h.Req.PostID}
	if err := database.GetInstance().Create(like).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
