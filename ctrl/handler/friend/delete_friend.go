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

func DeleteFriendHandle(c *gin.Context) {
	hd := &DeleteFriendHandler{}
	handler.Handle(c, hd)
}

type DeleteFriendHandler struct {
	Req  DeleteFriendRequest
	Resp DeleteFriendResponse
}

type DeleteFriendRequest struct {
	User     *database.User
	Friendid uint `json:"friendid"`
}

type DeleteFriendResponse struct {
	common.BaseResponse
}

func (h *DeleteFriendHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *DeleteFriendHandler) AfterBindReq() error {
	return nil
}

func (h *DeleteFriendHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *DeleteFriendHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *DeleteFriendHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *DeleteFriendHandler) NeedVerifyToken() bool {
	return true
}

func (h *DeleteFriendHandler) Process() {
	var friend *database.User
	// 判断是否为好友
	if err := database.GetInstance().Model(&h.Req.User).Association("Friends").Find(&friend, h.Req.Friendid); err != nil {
		msg := fmt.Sprintf("Judge relationship error, %v", err)
		seelog.Errorf("Judge relationship error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if friend.ID == 0 {
		seelog.Errorf("are not friends")
		h.SetError(common.ErrorInner, "are not friends")
		return
	}
	// 删除好友，更新数据库
	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(h.Req.User).Association("Friends").Delete(friend); err != nil {
		msg := fmt.Sprintf("delete friend error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(friend).Association("Friends").Delete(h.Req.User); err != nil {
		msg := fmt.Sprintf("delete friend error, %v", err)
		seelog.Errorf(friend.CreatedAt.Month().String())
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.User{}).Where("id = ?", friend.ID).Save(friend).Error; err != nil {
		msg := fmt.Sprintf("update friend error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.User.ID).Save(h.Req.User).Error; err != nil {
		msg := fmt.Sprintf("update user error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	tx.Commit()

}
