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

func AddMemberHandle(c *gin.Context) {
	hd := &AddMemberHandler{}
	handler.Handle(c, hd)
}

type AddMemberHandler struct {
	Req  AddMemberRequest
	Resp AddMemberResponse
}

type AddMemberRequest struct {
	User     *database.User
	Groupid  uint `json:"groupid"`
	Memberid uint `json:"memberid"`
}

type AddMemberResponse struct {
	common.BaseResponse
}

func (h *AddMemberHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *AddMemberHandler) AfterBindReq() error {
	return nil
}

func (h *AddMemberHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *AddMemberHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *AddMemberHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *AddMemberHandler) NeedVerifyToken() bool {
	return true
}

func (h *AddMemberHandler) Process() {
	// 查询群组
	var group *database.Group
	if err := database.GetInstance().Model(&database.Group{}).Preload("Members").Where("id = ?", h.Req.Groupid).Find(&group).Error; err != nil {
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
	if member.ID != 0 {
		seelog.Errorf("is member")
		h.SetError(common.ErrorInner, "is member")
		return
	}

	// 查询用户
	if err := database.GetInstance().Model(&database.User{}).Find(&member, h.Req.Memberid).Error; err != nil {
		msg := fmt.Sprintf("get user error, %v", err)
		seelog.Errorf("get user error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}

	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//添加用户与群组对应关系
	if err := tx.Model(group).Association("Members").Append(member); err != nil {
		msg := fmt.Sprintf("add group member error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(member).Association("Groups").Append(group); err != nil {
		msg := fmt.Sprintf("add group member error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	// // 保存数据库
	if err := tx.Model(&database.Group{}).Where("id = ?", h.Req.Groupid).Save(group).Error; err != nil {
		msg := fmt.Sprintf("save group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.User{}).Where("id = ?", h.Req.Memberid).Save(member).Error; err != nil {
		msg := fmt.Sprintf("save member error, %v", err)
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
