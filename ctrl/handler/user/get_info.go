package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/gin-gonic/gin"
)

func GetInfoHandle(c *gin.Context) {
	hd := &GetInfoHandler{}
	handler.Handle(c, hd)
}

type GetInfoHandler struct {
	Req  GetInfoRequest
	Resp GetInfoResponse
}

type GetInfoRequest struct {
	User *database.User
}

type GetInfoResponse struct {
	common.BaseResponse
	ID           uint   `json:"id"`
	Name         string `json:"username"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
	Follower     int64  `json:"follower"`
	Following    int64  `json:"following"`
}

func (h *GetInfoHandler) BindReq(c *gin.Context) error {
	return nil
}

func (h *GetInfoHandler) AfterBindReq() error {
	return nil
}

func (h *GetInfoHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetInfoHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetInfoHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetInfoHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetInfoHandler) Process() {
	h.Resp.ID = h.Req.User.ID
	h.Resp.Name = h.Req.User.Name
	h.Resp.Introduction = h.Req.User.Introduction
	h.Resp.Avatar = h.Req.User.Avatar
	if err := database.GetInstance().Model(&database.Follow{}).Where("user_id = ?", h.Req.User.ID).Count(&h.Resp.Following).Error; err != nil {
		msg := fmt.Sprintf("get following error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if err := database.GetInstance().Model(&database.Follow{}).Where("following_id = ?", h.Req.User.ID).Count(&h.Resp.Follower).Error; err != nil {
		msg := fmt.Sprintf("get follower error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
