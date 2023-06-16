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

func MemberQuitHandle(c *gin.Context) {
	hd := &MemberQuitHandler{}
	handler.Handle(c, hd)
}

type MemberQuitHandler struct {
	Req  MemberQuitRequest
	Resp MemberQuitResponse
}

type MemberQuitRequest struct {
	User    *database.User
	Groupid uint `json:"groupid"`
}

type MemberQuitResponse struct {
	common.BaseResponse
}

func (h *MemberQuitHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *MemberQuitHandler) AfterBindReq() error {
	return nil
}

func (h *MemberQuitHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *MemberQuitHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *MemberQuitHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *MemberQuitHandler) NeedVerifyToken() bool {
	return true
}

func (h *MemberQuitHandler) Process() {
	// 查询群组
	var group *database.Group
	if err := database.GetInstance().Model(&database.Group{}).Preload("Members").Find(&group, h.Req.Groupid).Error; err != nil {
		msg := fmt.Sprintf("get group error, %v", err)
		seelog.Errorf("get group error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	// 判断是否为成员
	var member *database.User
	if err := database.GetInstance().Model(&group).Association("Members").Find(&member, h.Req.User.ID); err != nil {
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
	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.User.ID).Save(member).Error; err != nil {
		msg := fmt.Sprintf("update member error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	tx.Commit()

}
