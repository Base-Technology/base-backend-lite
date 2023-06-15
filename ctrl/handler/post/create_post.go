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

func CreatePostHandle(c *gin.Context) {
	hd := &CreatePostHandler{}
	handler.Handle(c, hd)
}

type CreatePostHandler struct {
	Req  CreatePostRequest
	Resp CreatePostResponse
}

type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
	PostId  uint     `json:"post_id"`
	User    *database.User
}

type CreatePostResponse struct {
	common.BaseResponse
}

func (h *CreatePostHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *CreatePostHandler) AfterBindReq() error {
	return nil
}

func (h *CreatePostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *CreatePostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *CreatePostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *CreatePostHandler) NeedVerifyToken() bool {
	return true
}

func (h *CreatePostHandler) Process() {
	tx := database.GetInstance().Begin()
	defer func() {
		if h.Resp.BaseResponse.Code != 0 {
			if err := tx.Rollback().Error; err != nil {
				seelog.Errorf("tx rollback error, %v", err)
			}
		} else {
			if err := tx.Commit().Error; err != nil {
				msg := fmt.Sprintf("tx commit error, %v", err)
				seelog.Errorf(msg)
				h.SetError(common.ErrorInner, msg)
				return
			}
		}
	}()

	post := &database.Post{
		CreatorID:     h.Req.User.ID,
		Title:         h.Req.Title,
		Content:       h.Req.Content,
		PostPointedID: h.Req.PostId,
	}
	fields := []string{"creator_id", "title", "content", "create_at", "update_at"}
	if post.PostPointedID != 0 {
		fields = append(fields, "post_pointed_id")
	}
	if err := tx.Select("id", fields).Create(&post).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	for _, image := range h.Req.Images {
		if err := tx.Create(&database.Image{
			PostID: post.ID,
			Source: image,
		}).Error; err != nil {
			msg := fmt.Sprintf("insert to database error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
	}
}
