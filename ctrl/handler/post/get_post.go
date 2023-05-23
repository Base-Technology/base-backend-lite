package post

import (
	"fmt"
	"time"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/gin-gonic/gin"
)

func GetPostHandle(c *gin.Context) {
	hd := &GetPostHandler{}
	handler.Handle(c, hd)
}

type GetPostHandler struct {
	Req  GetPostRequest
	Resp GetPostResponse
}

type GetPostRequest struct {
	Type string
	User *database.User
}

type GetPostResponse struct {
	common.BaseResponse
	Posts []*PostDeatail `json:"data"`
}

type PostDeatail struct {
	ID       uint      `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
}

func (h *GetPostHandler) BindReq(c *gin.Context) error {
	h.Req.Type = c.Query("type")
	return nil
}

func (h *GetPostHandler) AfterBindReq() error {
	return nil
}

func (h *GetPostHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetPostHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetPostHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetPostHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetPostHandler) Process() {
	var err error
	posts := []*database.Post{}
	switch h.Req.Type {
	case "me":
		err = database.GetInstance().Model(&database.Post{}).Where("creator_id = ?", h.Req.User.ID).Find(&posts).Error
	case "like":
		likes := []*database.Like{}
		err = database.GetInstance().Preload("Post").Where("user_id = ?", h.Req.User.ID).Find(&likes).Error
		for _, like := range likes {
			posts = append(posts, like.Post)
		}
	case "collect":
		collects := []*database.Collect{}
		err = database.GetInstance().Preload("Post").Where("user_id = ?", h.Req.User.ID).Find(&collects).Error
		for _, collect := range collects {
			posts = append(posts, collect.Post)
		}
	default:
		msg := fmt.Sprintf("invalid type: [%v]", h.Req.Type)
		h.SetError(common.ErrorInvalidParams, msg)
		return
	}

	if err != nil {
		msg := fmt.Sprintf("get post error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Posts = []*PostDeatail{}
	for _, post := range posts {
		h.Resp.Posts = append(h.Resp.Posts, &PostDeatail{
			ID:       post.ID,
			Title:    post.Title,
			Content:  post.Content,
			CreateAt: post.CreatedAt,
		})
	}
}
