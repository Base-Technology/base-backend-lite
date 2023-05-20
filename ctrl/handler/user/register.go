package user

import (
	"fmt"

	"github.com/Base-Technology/base-app-lite/common"
	"github.com/Base-Technology/base-app-lite/ctrl/handler"
	"github.com/Base-Technology/base-app-lite/database"
	"github.com/Base-Technology/base-app-lite/seelog"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandle(c *gin.Context) {
	hd := &RegisterHandler{}
	handler.Handle(c, hd)
}

type RegisterHandler struct {
	Req  RegisterRequest
	Resp RegisterResponse
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Area     string `json:"area" binding:"required"`
	School   string `json:"school" binding:"required"`
}

type RegisterResponse struct {
	common.BaseResponse
}

func (h *RegisterHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *RegisterHandler) AfterBindReq() error {
	return nil
}

func (h *RegisterHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *RegisterHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *RegisterHandler) SetUser(user *database.User) {
}

func (h *RegisterHandler) NeedVerifyToken() bool {
	return false
}

func (h *RegisterHandler) Process() {
	// hash the password
	hp, err := bcrypt.GenerateFromPassword([]byte(h.Req.Password), bcrypt.DefaultCost)
	if err != nil {
		msg := fmt.Sprintf("hash password error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// generate private key
	k, err := crypto.GenerateKey()
	if err != nil {
		msg := fmt.Sprintf("generate private key error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// insert into database
	user := &database.User{
		Phone:      h.Req.Phone,
		Name:       h.Req.Username,
		Password:   string(hp),
		Area:       h.Req.Area,
		School:     h.Req.School,
		PrivateKey: hexutil.Encode(crypto.FromECDSA(k)),
	}
	if err := database.GetInstance().Create(user).Error; err != nil {
		msg := fmt.Sprintf("insert into database error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
