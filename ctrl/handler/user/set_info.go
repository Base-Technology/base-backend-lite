package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func SetInfoHandle(c *gin.Context) {
	hd := &SetInfoHandler{}
	handler.Handle(c, hd)
}

type SetInfoHandler struct {
	Req  SetInfoRequest
	Resp SetInfoResponse
}

type SetInfoRequest struct {
	Username     string `json:"username" binding:"required"`
	Introduction string `json:"introduction"`
	Avatar       string `json:"avatar"`
	User         *database.User
}

type SetInfoResponse struct {
	common.BaseResponse
}

func (h *SetInfoHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *SetInfoHandler) AfterBindReq() error {
	return nil
}

func (h *SetInfoHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *SetInfoHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *SetInfoHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *SetInfoHandler) NeedVerifyToken() bool {
	return true
}

func (h *SetInfoHandler) Process() {
	if err := database.GetInstance().Model(h.Req.User).Updates(map[string]interface{}{
		"Name":         h.Req.Username,
		"Introduction": h.Req.Introduction,
		"Avatar":       h.Req.Avatar,
	}).Error; err != nil {
		msg := fmt.Sprintf("update to database error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
