package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/imtp"
	"github.com/Base-Technology/base-backend-lite/school"
	"github.com/Base-Technology/base-backend-lite/seelog"
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
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Area         string `json:"area" binding:"required"`
	School       string `json:"school" binding:"required"`
	ValidateCode string `json:"validate_code" binding:"required"`
	Avatar       string `json:"avatar"`
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
	// check the validate code
	code, ok := validateCodes.Get(h.Req.Phone)
	if !ok || code != h.Req.ValidateCode {
		msg := fmt.Sprintf("validate code invalid")
		h.SetError(common.ErrorInvalidParams, msg)
		return
	}
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
	kBytes := hexutil.Encode(crypto.FromECDSA(k))
	// login to create the account
	_, _, err = imtp.Login(kBytes[2:])
	if err != nil {
		msg := fmt.Sprintf("login imtp error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// invite to imtp group
	address := crypto.PubkeyToAddress(k.PublicKey).Hex()
	if err := school.InviteUserToSchoolGroup(address, h.Req.School); err != nil {
		msg := fmt.Sprintf("invite user to imtp group error, %v", err)
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
		PrivateKey: kBytes,
		IMTPUserID: imtp.GetUserIDFromAddress(address),
		Avatar:     h.Req.Avatar,
	}
	if err := database.GetInstance().Create(user).Error; err != nil {
		msg := fmt.Sprintf("insert into database error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	validateCodes.Delete(h.Req.Phone)
}
