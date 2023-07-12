package post

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
	Comment []*types.CommentDetail `json:"data"`
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
	if err := database.GetInstance().Model(&database.Comment{}).Preload("Creator").Where("post_id = ?", h.Req.PostId).Order("created_at desc").Offset((h.Req.Page - 1) * h.Req.Limit).Limit(h.Req.Limit).Find(&comments).Error; err != nil {
		msg := fmt.Sprintf("get comments error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	for _, c := range comments {
		var pointed uint
		if c.CommentPointedID == nil {
			pointed = 0
		} else {
			pointed = *c.CommentPointedID
		}
		h.Resp.Comment = append(h.Resp.Comment, &types.CommentDetail{
			CommentID:        c.ID,
			CreatorID:        c.CreatorID,
			CreatorName:      c.Creator.Name,
			CreatorAvatar:    c.Creator.Avatar,
			CreateAt:         c.CreatedAt,
			Content:          c.Content,
			Level:            c.Level,
			CommentPointedID: pointed,
		})
	}
}
