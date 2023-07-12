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

func FollowHandle(c *gin.Context) {
	hd := &FollowHandler{}
	handler.Handle(c, hd)
}

type FollowHandler struct {
	Req  FollowRequest
	Resp FollowResponse
}

type FollowRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	User   *database.User
}

type FollowResponse struct {
	common.BaseResponse
}

func (h *FollowHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *FollowHandler) AfterBindReq() error {
	return nil
}

func (h *FollowHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *FollowHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *FollowHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *FollowHandler) NeedVerifyToken() bool {
	return true
}

func (h *FollowHandler) Process() {
	follow := &database.Follow{UserID: h.Req.User.ID, FollowingID: h.Req.UserID}
	if err := database.GetInstance().Create(follow).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
