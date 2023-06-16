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

func CreateGroupHandle(c *gin.Context) {
	hd := &CreateGroupHandler{}
	handler.Handle(c, hd)
}

type CreateGroupHandler struct {
	Req  CreateGroupRequest
	Resp CreateGroupResponse
}

type CreateGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description "`
	MembersID   []uint `json:"membersID"`
	Avatar      string `json:"avatat"`
	User        *database.User
}

type CreateGroupResponse struct {
	common.BaseResponse
	Groupid uint `json:"data"`
}

func (h *CreateGroupHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *CreateGroupHandler) AfterBindReq() error {
	return nil
}

func (h *CreateGroupHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *CreateGroupHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *CreateGroupHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *CreateGroupHandler) NeedVerifyToken() bool {
	return true
}

func (h *CreateGroupHandler) Process() {
	var user *database.User
	var members []*database.User
	// 设置回退
	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建群组，插入数据库
	group := &database.Group{Name: h.Req.Name, Description: h.Req.Description, Members: members, CreatorID: h.Req.User.ID, Avatar: h.Req.Avatar}

	if err := tx.Create(group).Error; err != nil {
		msg := fmt.Sprintf("insert Group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	// 将群主添加到群组中
	if err := tx.Model(group).Association("Members").Append(h.Req.User); err != nil {
		msg := fmt.Sprintf("add member to group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	// 将群组添加到用户中
	if err := tx.Model(h.Req.User).Association("Groups").Append(group); err != nil {
		msg := fmt.Sprintf("add group to user error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	// 获取成员并存入数据库
	for _, userID := range h.Req.MembersID {
		if err := tx.Model(&database.User{}).Where("id = ? ", userID).Find(&user); err != nil {
			msg := fmt.Sprintf("get user error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			return
		}
		if err := tx.Model(user).Association("Groups").Append(group); err != nil {
			msg := fmt.Sprintf("add group to user error, %v", err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			tx.Rollback()
			return
		}
		members = append(members, user)
	}

	// 保存到数据库
	if err := tx.Model(&database.Group{}).Where("id = ? ", group.ID).Save(group).Error; err != nil {
		msg := fmt.Sprintf("save Group error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	if err := tx.Model(&database.User{}).Where("id = ? ", h.Req.User.ID).Save(h.Req.User).Error; err != nil {
		msg := fmt.Sprintf("save User error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	tx.Commit()
	h.Resp.Groupid = group.ID
}
