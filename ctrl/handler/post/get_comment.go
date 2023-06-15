package post

import (
	"fmt"
	"time"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetCommentHandle(c *gin.Context) {
	hd := &GetCommentHandler{}
	handler.Handle(c, hd)
}

type GetCommentHandler struct {
	Req  GetCommentRequest
	Resp GetCommentResponse
}

type GetCommentRequest struct {
	PostId uint `json:"post_id"`
	User   *database.User
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}

type GetCommentResponse struct {
	common.BaseResponse
	Comment []*CommentDeatail `json:"data"`
}
type CommentDeatail struct {
	CommentID        uint      `json:"comment_id"`
	CreatorID        uint      `json:"creator_id"`
	CreatorName      string    `json:"creator_name"`
	CreatorAvatar    string    `json:"creator_avatar"`
	CreateAt         time.Time `json:"create_at"`
	Content          string    `json:"content"`
	Level            uint      `json:"level"`
	CommentPointedID uint      `json:"commentpointed_id"`
}

func (h *GetCommentHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *GetCommentHandler) AfterBindReq() error {
	return nil
}

func (h *GetCommentHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetCommentHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetCommentHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetCommentHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetCommentHandler) Process() {
	var comments []*database.Comment
	if err := database.GetInstance().Model(&database.Comment{}).Where("post_id = ?", h.Req.PostId).Offset((h.Req.Page - 1) * h.Req.Limit).Limit(h.Req.Limit).Find(&comments).Error; err != nil {
		msg := fmt.Sprintf("get comments error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	for _, c := range comments {
		var creator *database.User
		if err := database.GetInstance().Model(&database.User{}).Where("id = ?", c.CreatorID).Find(&creator).Error; err != nil {
			msg := fmt.Sprintf("get comment creator error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		var pointed uint
		if c.CommentPointedID == nil {
			pointed = 0
		} else {
			pointed = *c.CommentPointedID
		}
		h.Resp.Comment = append(h.Resp.Comment, &CommentDeatail{
			CommentID:        c.ID,
			CreatorID:        c.CreatorID,
			CreatorName:      creator.Name,
			CreatorAvatar:    creator.Avatar,
			CreateAt:         c.CreatedAt,
			Content:          c.Content,
			Level:            c.Level,
			CommentPointedID: pointed,
		})
	}
}
