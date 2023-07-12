package friend

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/ctrl/types"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
)

func GetFriendListHandle(c *gin.Context) {
	hd := &GetFriendListHandler{}
	handler.Handle(c, hd)
}

type GetFriendListHandler struct {
	Req  GetFriendListRequest
	Resp GetFriendListResponse
}

type GetFriendListRequest struct {
	User *database.User
}

type GetFriendListResponse struct {
	common.BaseResponse
	FriendList []*types.UserDetail `json:"data"`
}

func (h *GetFriendListHandler) BindReq(c *gin.Context) error {
	return nil
}

func (h *GetFriendListHandler) AfterBindReq() error {
	return nil
}

func (h *GetFriendListHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetFriendListHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetFriendListHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetFriendListHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetFriendListHandler) Process() {
	var err error
	var user *database.User
	if err = database.GetInstance().Model(&database.User{}).Preload("Friends").Where("id = ?", h.Req.User.ID).Find(&user).Error; err != nil {
		msg := fmt.Sprintf("get friend list error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}

	for _, u := range user.Friends {
		h.Resp.FriendList = append(h.Resp.FriendList, &types.UserDetail{
			ID:     u.ID,
			Name:   u.Name,
			Avatar: u.Avatar,
			Sex:    u.Sex,
		})
	}
}
