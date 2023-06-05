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

func RequestAddFriendHandle(c *gin.Context) {
	hd := &RequestAddFriendHandler{}
	handler.Handle(c, hd)
}

type RequestAddFriendHandler struct {
	Req  RequestAddFriendRequest
	Resp RequestAddFriendResponse
}

type RequestAddFriendRequest struct {
	User     *database.User
	Friendid uint   `json:"friendid"`
	Message  string `json:"message"`
}

type RequestAddFriendResponse struct {
	common.BaseResponse
}

func (h *RequestAddFriendHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *RequestAddFriendHandler) AfterBindReq() error {
	return nil
}

func (h *RequestAddFriendHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *RequestAddFriendHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *RequestAddFriendHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *RequestAddFriendHandler) NeedVerifyToken() bool {
	return true
}

func (h *RequestAddFriendHandler) Process() {
	var err error
	var user *database.User
	// 判断是否为好友
	if err := database.GetInstance().Model(&h.Req.User).Association("Friends").Find(&user, h.Req.Friendid); err != nil {
		msg := fmt.Sprintf("Judge relationship error, %v", err)
		seelog.Errorf("Judge relationship error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if user.ID != 0 {
		seelog.Errorf("are already friend")
		h.SetError(common.ErrorInner, "are already friend")
		return
	}
	// 查询是否有好友请求，如果有修改请求，更新数据库
	var request *database.FriendRequest
	if err = database.GetInstance().Model(&database.FriendRequest{}).Where("user_id = ? AND sender_ID = ?", h.Req.Friendid, h.Req.User.ID).Find(&request).Error; err != nil {
		msg := fmt.Sprintf("get request error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	if request.UserID != 0 {
		if request.Status == "pending" {
			seelog.Errorf("request sent")
			h.SetError(common.ErrorInner, "request sent")
			return
		}
		request.Status = "pending"
		request.Message = h.Req.Message
		if err := database.GetInstance().Model(&database.FriendRequest{}).Where("user_id = ? AND sender_ID = ?", h.Req.Friendid, h.Req.User.ID).Save(request).Error; err != nil {
			msg := fmt.Sprintf("update friendRequest error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		return
	}
	// 添加好友请求，插入数据库
	addrequest := &database.FriendRequest{UserID: h.Req.Friendid, SenderID: h.Req.User.ID, Message: h.Req.Message, Name: h.Req.User.Name, Avatar: h.Req.User.Avatar, Status: "pending"}
	if err := database.GetInstance().Create(addrequest).Error; err != nil {
		msg := fmt.Sprintf("insert to database error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
}
