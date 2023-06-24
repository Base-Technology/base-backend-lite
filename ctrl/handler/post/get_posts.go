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

func GetPostHandle(c *gin.Context) {
	hd := &GetPostHandler{}
	handler.Handle(c, hd)
}

type GetPostHandler struct {
	Req  GetPostRequest
	Resp GetPostResponse
}

type GetPostRequest struct {
	Type  string
	Page  int `json:"page"`
	Limit int `json:"limit"`
	User  *database.User
}

type GetPostResponse struct {
	common.BaseResponse
	Posts []*types.PostDetail `json:"data"`
}

func (h *GetPostHandler) BindReq(c *gin.Context) error {
	h.Req.Type = c.Query("type")
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
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
		err = database.GetInstance().Model(&database.Post{}).Preload("Creator").Where("creator_id = ?", h.Req.User.ID).Order("created_at desc").Offset((h.Req.Page - 1) * h.Req.Limit).Limit(h.Req.Limit).Find(&posts).Error
	case "like":
		var post_id []uint
		likes := []*database.Like{}
		err := database.GetInstance().Model(&database.Like{}).Select("post_id").Where("user_id = ?", h.Req.User.ID).Find(&post_id).Error
		if err != nil {
			msg := fmt.Sprintf("get post id error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		err = database.GetInstance().Model(&database.Post{}).Preload("Creator").Where("id IN (?)", post_id).
			Order("created_at desc").
			Offset((h.Req.Page - 1) * h.Req.Limit).
			Limit(h.Req.Limit).
			Find(&posts).Error
		if err != nil {
			msg := fmt.Sprintf("get post error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		for _, like := range likes {
			posts = append(posts, like.Post)
		}
	case "collect":
		var post_id []uint
		collects := []*database.Collect{}
		err := database.GetInstance().Model(&database.Collect{}).Select("post_id").Where("user_id = ?", h.Req.User.ID).Find(&post_id).Error
		if err != nil {
			msg := fmt.Sprintf("get post id error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		err = database.GetInstance().Model(&database.Post{}).Preload("Creator").Where("id IN (?)", post_id).
			Order("created_at desc").
			Offset((h.Req.Page - 1) * h.Req.Limit).
			Limit(h.Req.Limit).
			Find(&posts).Error
		if err != nil {
			msg := fmt.Sprintf("get post error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		for _, collect := range collects {
			posts = append(posts, collect.Post)
		}
	case "all":
		err = database.GetInstance().Model(&database.Post{}).Preload("Creator").Order("created_at desc").Offset((h.Req.Page - 1) * h.Req.Limit).Limit(h.Req.Limit).Find(&posts).Error
	case "following":
		var followings []uint
		if err := database.GetInstance().Model(&database.Follow{}).Select("following_id").Where("user_id = ?", h.Req.User.ID).Find(&followings).Error; err != nil {
			msg := fmt.Sprintf("get followings error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		err = database.GetInstance().Model(&database.Post{}).Preload("Creator").Where("creator_id IN (?)", followings).
			Order("created_at desc").
			Offset((h.Req.Page - 1) * h.Req.Limit).
			Limit(h.Req.Limit).
			Find(&posts).Error
	default:
		msg := fmt.Sprintf("invalid type: [%v]", h.Req.Type)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	if err != nil {
		msg := fmt.Sprintf("get post error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Posts = []*types.PostDetail{}
	for _, post := range posts {
		var comment_count int64
		var like_count int64
		var collect_count int64
		if err := database.GetInstance().Model(&database.Comment{}).Where("post_id = ? AND level = ?", post.ID, 1).Count(&comment_count).Error; err != nil {
			msg := fmt.Sprintf("get comment count error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		if err := database.GetInstance().Model(&database.Like{}).Where("post_id = ?", post.ID).Count(&like_count).Error; err != nil {
			msg := fmt.Sprintf("get like count error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		if err := database.GetInstance().Model(&database.Collect{}).Where("post_id = ?", post.ID).Count(&collect_count).Error; err != nil {
			msg := fmt.Sprintf("get collect count error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		h.Resp.Posts = append(h.Resp.Posts, &types.PostDetail{
			ID:            post.ID,
			Title:         post.Title,
			Content:       post.Content,
			CreateAt:      post.CreatedAt,
			CreatorId:     post.CreatorID,
			CreatorName:   post.Creator.Name,
			CreatorAvatar: post.Creator.Avatar,
			CommentCount:  comment_count,
			LikeCount:     like_count,
			CollectCount:  collect_count,
		})
	}
}
