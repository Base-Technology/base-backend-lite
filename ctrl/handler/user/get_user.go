package user

import (
	"fmt"

	"github.com/Base-Technology/base-app-lite/common"
	"github.com/Base-Technology/base-app-lite/ctrl/handler"
	"github.com/Base-Technology/base-app-lite/database"
	"github.com/gin-gonic/gin"
)

func GetUserHandle(c *gin.Context) {
	hd := &GetUserHandler{}
	handler.Handle(c, hd)
}

type GetUserHandler struct {
	Req  GetUserRequest
	Resp GetUserResponse
}

type GetUserRequest struct {
	Type string
	User *database.User
}

type GetUserResponse struct {
	common.BaseResponse
	Users []*UserDeatail `json:"data"`
}

type UserDeatail struct {
	ID           uint   `json:"id"`
	Name         string `json:"username"`
	Area         string `json:"area"`
	School       string `json:"school"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
}

func (h *GetUserHandler) BindReq(c *gin.Context) error {
	h.Req.Type = c.Query("type")
	return nil
}

func (h *GetUserHandler) AfterBindReq() error {
	return nil
}

func (h *GetUserHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetUserHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetUserHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetUserHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetUserHandler) Process() {
	var err error
	users := []*database.User{}
	switch h.Req.Type {
	case "following":
		follows := []*database.Follow{}
		err = database.GetInstance().Preload("Following").Where("user_id = ?", h.Req.User.ID).Find(&follows).Error
		for _, follower := range follows {
			users = append(users, follower.Following)
		}
	case "follower":
		follows := []*database.Follow{}
		err = database.GetInstance().Preload("User").Where("following_id = ?", h.Req.User.ID).Find(&follows).Error
		for _, follower := range follows {
			users = append(users, follower.User)
		}
	default:
		msg := fmt.Sprintf("invalid type: [%v]", h.Req.Type)
		h.SetError(common.ErrorInvalidParams, msg)
		return
	}

	if err != nil {
		msg := fmt.Sprintf("get user error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Users = []*UserDeatail{}
	for _, user := range users {
		h.Resp.Users = append(h.Resp.Users, &UserDeatail{
			ID:           user.ID,
			Name:         user.Name,
			Area:         user.Area,
			School:       user.School,
			Introduction: user.Introduction,
			Avatar:       user.Avatar,
		})
	}
}
