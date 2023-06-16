package group

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func RemoveMemberHandle(c *gin.Context) {
	hd := &RemoveMemberHandler{}
	handler.Handle(c, hd)
}

type RemoveMemberHandler struct {
	Req  RemoveMemberRequest
	Resp RemoveMemberResponse
}

type RemoveMemberRequest struct {
	User     *database.User
	Groupid  uint `json:"groupid"`
	Memberid uint `json:"memberid"`
}

type RemoveMemberResponse struct {
	common.BaseResponse
}

func (h *RemoveMemberHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *RemoveMemberHandler) AfterBindReq() error {
	return nil
}

func (h *RemoveMemberHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *RemoveMemberHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *RemoveMemberHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *RemoveMemberHandler) NeedVerifyToken() bool {
	return true
}

func (h *RemoveMemberHandler) Process() {
	// 查询群组
	var group *database.Group
	if err := database.GetInstance().Model(&database.Group{}).Preload("Members").Find(&group, h.Req.Groupid).Error; err != nil {
		msg := fmt.Sprintf("get group error, %v", err)
		seelog.Errorf("get group error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// 判断是否为群主
	if group.CreatorID != h.Req.User.ID {
		seelog.Errorf("no group owner")
		h.SetError(common.ErrorInner, "no group owner")
		return
	}
	// 判断是否为成员
	var member *database.User
	if err := database.GetInstance().Model(&group).Association("Members").Find(&member, h.Req.Memberid); err != nil {
		msg := fmt.Sprintf("get groupmenber error, %v", err)
		seelog.Errorf("get groupmenber error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	if member.ID == 0 {
		seelog.Errorf("no member")
		h.SetError(common.ErrorInner, "no member")
		return
	}
	// 删除成员，更新数据库
	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Model(group).Association("Members").Delete(member); err != nil {
		msg := fmt.Sprintf("delete group`s member error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(member).Association("Groups").Delete(group); err != nil {
		msg := fmt.Sprintf("delete member`s group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.Group{}).Where("id = ?", h.Req.Groupid).Save(group).Error; err != nil {
		msg := fmt.Sprintf("update group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.Memberid).Save(member).Error; err != nil {
		msg := fmt.Sprintf("update member error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	tx.Commit()

}
