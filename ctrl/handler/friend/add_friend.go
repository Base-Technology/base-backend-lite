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

func AddFriendHandle(c *gin.Context) {
	hd := &AddFriendHandler{}
	handler.Handle(c, hd)
}

type AddFriendHandler struct {
	Req  AddFriendRequest
	Resp AddFriendResponse
}

type AddFriendRequest struct {
	User     *database.User
	Senderid uint `json:"senderid"`
}

type AddFriendResponse struct {
	common.BaseResponse
}

func (h *AddFriendHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *AddFriendHandler) AfterBindReq() error {
	return nil
}

func (h *AddFriendHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *AddFriendHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *AddFriendHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *AddFriendHandler) NeedVerifyToken() bool {
	return true
}

func (h *AddFriendHandler) Process() {
	var err error
	// 获取好友请求，判断请求是否待处理
	var request *database.FriendRequest
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
	var sender *database.User
	// 添加好友，更新数据库User表
	if err = database.GetInstance().Model(&database.User{}).Where("id = ?", h.Req.Senderid).Find(&sender).Error; err != nil {
		msg := fmt.Sprintf("get user info error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(sender).Association("Friends").Append(h.Req.User); err != nil {
		msg := fmt.Sprintf("update sender`s friend error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	if err := tx.Model(h.Req.User).Association("Friends").Append(sender); err != nil {
		msg := fmt.Sprintf("update user`s friend error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.Senderid).Save(sender).Error; err != nil {
		msg := fmt.Sprintf("save sender error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.User.ID).Save(h.Req.User).Error; err != nil {
		msg := fmt.Sprintf("save user error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	request.Status = "accepted"
	if err := tx.Model(&database.FriendRequest{}).Where("user_id = ? AND sender_ID = ?", h.Req.User.ID, h.Req.Senderid).Save(request).Error; err != nil {
		msg := fmt.Sprintf("update request error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	if err := tx.Commit().Error; err != nil {
		msg := fmt.Sprintf("update database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

}
