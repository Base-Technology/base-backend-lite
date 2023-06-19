package user

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/detail"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/gin-gonic/gin"
)

func GetOtherUserHandle(c *gin.Context) {
	hd := &GetOtherUserHandler{}
	handler.Handle(c, hd)
}

type GetOtherUserHandler struct {
	Req  GetOtherUserRequest
	Resp GetOtherUserResponse
}

type GetOtherUserRequest struct {
	IMTPUserID string
	User       *database.User
}

type GetOtherUserResponse struct {
	common.BaseResponse
	User detail.UserDeatailMore `json:"user"`
}

func (h *GetOtherUserHandler) BindReq(c *gin.Context) error {
	h.Req.IMTPUserID = c.Query("imtp_user_id")
	return nil
}

func (h *GetOtherUserHandler) AfterBindReq() error {
	return nil
}

func (h *GetOtherUserHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetOtherUserHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetOtherUserHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetOtherUserHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetOtherUserHandler) Process() {
	users := []*database.User{}
	if err := database.GetInstance().Model(&database.User{}).Where("imtp_user_id = ?", h.Req.IMTPUserID).Find(&users).Error; err != nil {
		msg := fmt.Sprintf("get user error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if len(users) == 0 {
		msg := fmt.Sprintf("user not found")
		h.SetError(common.ErrorInner, msg)
		return
	}
	user := users[0]
	h.Resp.User = detail.UserDeatailMore{
		ID:           user.ID,
		Name:         user.Name,
		Area:         user.Area,
		School:       user.School,
		Introduction: user.Introduction,
		Avatar:       user.Avatar,
		IMTPUserID:   user.IMTPUserID,
		Sex:          user.Sex,
	}
}
