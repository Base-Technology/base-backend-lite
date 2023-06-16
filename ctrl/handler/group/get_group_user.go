package group

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/gin-gonic/gin"
)

func GetGroupUserHandle(c *gin.Context) {
	hd := &GetGroupUserHandler{}
	handler.Handle(c, hd)
}

type GetGroupUserHandler struct {
	Req  GetGroupUserRequest
	Resp GetGroupUserResponse
}

type GetGroupUserRequest struct {
	Type string
	User *database.User
}

type GetGroupUserResponse struct {
	common.BaseResponse
	Users []*UserDeatail `json:"data"`
}

type UserDeatail struct {
	ID     uint   `json:"id"`
	Name   string `json:"username"`
	Avatar string `json:"avatar"`
}

func (h *GetGroupUserHandler) BindReq(c *gin.Context) error {
	h.Req.Type = c.Query("type")
	return nil
}

func (h *GetGroupUserHandler) AfterBindReq() error {
	return nil
}

func (h *GetGroupUserHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetGroupUserHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetGroupUserHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetGroupUserHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetGroupUserHandler) Process() {
	users := []*database.User{}
	if err := database.GetInstance().Model(&database.User{}).Where("school = ?", h.Req.User.School).Find(&users).Error; err != nil {
		msg := fmt.Sprintf("get group user error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Users = []*UserDeatail{}
	for _, user := range users {
		h.Resp.Users = append(h.Resp.Users, &UserDeatail{
			ID:     user.ID,
			Name:   user.Name,
			Avatar: user.Avatar,
		})
	}
}
