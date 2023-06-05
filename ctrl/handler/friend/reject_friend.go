package friend

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func RejectFriendHandle(c *gin.Context) {
	hd := &RejectFriendHandler{}
	handler.Handle(c, hd)
}

type RejectFriendHandler struct {
	Req  RejectFriendRequest
	Resp RejectFriendResponse
}

type RejectFriendRequest struct {
	User     *database.User
	Senderid uint `json:"senderid"`
}

type RejectFriendResponse struct {
	common.BaseResponse
}

func (h *RejectFriendHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *RejectFriendHandler) AfterBindReq() error {
	return nil
}

func (h *RejectFriendHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *RejectFriendHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *RejectFriendHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *RejectFriendHandler) NeedVerifyToken() bool {
	return true
}

func (h *RejectFriendHandler) Process() {
	var err error
	var request *database.FriendRequest
	// 获取好友请求，判断请求状态是否待处理
	if err = database.GetInstance().Model(&database.FriendRequest{}).Where("user_id = ? AND sender_ID = ?", h.Req.User.ID, h.Req.Senderid).Find(&request).Error; err != nil {
		msg := fmt.Sprintf("get request error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if request.UserID == 0 {
		seelog.Errorf("not request")
		h.SetError(common.ErrorInner, "not request")
		return
	}
	if request.Status != "pending" {
		seelog.Errorf("Request processed")
		h.SetError(common.ErrorInner, "Request processed")
		return
	}
	// 处理请求数据，修改数据库
	request.Status = "declined"
	if err := database.GetInstance().Model(&database.FriendRequest{}).Where("user_id = ? AND sender_ID = ?", h.Req.User.ID, h.Req.Senderid).Save(request).Error; err != nil {
		msg := fmt.Sprintf("update FriendRequest error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
