package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

func ResetPasswordHandle(c *gin.Context) {
	hd := &ResetPasswordHandler{}
	handler.Handle(c, hd)
}

type ResetPasswordHandler struct {
	Req  ResetPasswordRequest
	Resp ResetPasswordResponse
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
	User     *database.User
}

type ResetPasswordResponse struct {
	common.BaseResponse
}

func (h *ResetPasswordHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *ResetPasswordHandler) AfterBindReq() error {
	return nil
}

func (h *ResetPasswordHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *ResetPasswordHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *ResetPasswordHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *ResetPasswordHandler) NeedVerifyToken() bool {
	return true
}

func (h *ResetPasswordHandler) Process() {
	// hash the password
	hp, err := bcrypt.GenerateFromPassword([]byte(h.Req.Password), bcrypt.DefaultCost)
	if err != nil {
		msg := fmt.Sprintf("hash password error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// update to database
	if err := database.GetInstance().Model(h.Req.User).Update("Password", string(hp)).Error; err != nil {
		msg := fmt.Sprintf("update to database error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
