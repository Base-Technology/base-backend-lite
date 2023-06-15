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

func AddCommentHandle(c *gin.Context) {
	hd := &AddCommentHandler{}
	handler.Handle(c, hd)
}

type AddCommentHandler struct {
	Req  AddCommentRequest
	Resp AddCommentResponse
}

type AddCommentRequest struct {
	PostId           uint   `json:"post_id"`
	CommentPointedID uint   `json:"comment_pointed"`
	Content          string `json:"content"`
	User             *database.User
}

type AddCommentResponse struct {
	common.BaseResponse
	CommentID uint `json:"data"`
}

func (h *AddCommentHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *AddCommentHandler) AfterBindReq() error {
	return nil
}

func (h *AddCommentHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *AddCommentHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *AddCommentHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *AddCommentHandler) NeedVerifyToken() bool {
	return true
}

func (h *AddCommentHandler) Process() {
	var comment *database.Comment

	if h.Req.CommentPointedID == 0 {
		comment = &database.Comment{
			PostID:           h.Req.PostId,
			CreatorID:        h.Req.User.ID,
			CommentPointedID: nil,
			Content:          h.Req.Content,
			Level:            1,
		}
	} else {
		var commentpointed *database.Comment
		if err := database.GetInstance().Model(&database.Comment{}).Where("id = ?", h.Req.CommentPointedID).Find(&commentpointed).Error; err != err {
			msg := fmt.Sprintf("get comment pointed error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		if commentpointed.PostID != h.Req.PostId {
			seelog.Errorf("is not post comment")
			h.SetError(common.ErrorInner, "is not post comment")
			return
		}
		comment = &database.Comment{
			PostID:           h.Req.PostId,
			CreatorID:        h.Req.User.ID,
			CommentPointedID: &h.Req.CommentPointedID,
			Content:          h.Req.Content,
			Level:            commentpointed.Level + 1,
		}
	}
	var post *database.Post
	if err := database.GetInstance().Model(&database.Post{}).Where("id = ?", h.Req.PostId).Find(&post).Error; err != nil {
		msg := fmt.Sprintf("get post error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if err := database.GetInstance().Create(comment).Error; err != nil {
		msg := fmt.Sprintf("insert comment error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	h.Resp.CommentID = comment.ID
}
