package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/school"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/Base-Technology/base-backend-lite/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandle(c *gin.Context) {
	hd := &LoginHandler{}
	handler.Handle(c, hd)
}

type LoginHandler struct {
	Req  LoginRequest
	Resp LoginResponse
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	common.BaseResponse
	Token      string `json:"token"`
	UserID     int    `json:"user_id"`
	PrivateKey string `json:"private_key"`
	GroupID    string `json:"group_id"`
}

func (h *LoginHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *LoginHandler) AfterBindReq() error {
	return nil
}

func (h *LoginHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *LoginHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *LoginHandler) SetUser(user *database.User) {
}

func (h *LoginHandler) NeedVerifyToken() bool {
	return false
}

func (h *LoginHandler) Process() {
	user := &database.User{}
	if err := database.GetInstance().Where("phone = ?", h.Req.Phone).First(user).Error; err != nil {
		msg := fmt.Sprintf("get user error, %v", err)
		h.SetError(common.ErrorPassword, msg)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(h.Req.Password)); err != nil {
		h.SetError(common.ErrorPassword, "invalid password")
		return
	}

	t, err := token.GenerateToken(user.ID)
	if err != nil {
		msg := fmt.Sprintf("generate token error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	groupID, err := school.GetGroupIDByName(user.School)
	if err != nil {
		msg := fmt.Sprintf("get group id error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	h.Resp.Token = t
	h.Resp.UserID = int(user.ID)
	h.Resp.PrivateKey = user.PrivateKey
	h.Resp.GroupID = groupID
}
