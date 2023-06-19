package friend

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/detail"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
)

func GetFriendRequestListHandle(c *gin.Context) {
	hd := &GetFriendRequestListHandler{}
	handler.Handle(c, hd)
}

type GetFriendRequestListHandler struct {
	Req  GetFriendRequestListRequest
	Resp GetFriendRequestListResponse
}

type GetFriendRequestListRequest struct {
	User *database.User
}

type GetFriendRequestListResponse struct {
	common.BaseResponse
	FriendRequestList []*detail.RequestDetail `json:"data"`
}

func (h *GetFriendRequestListHandler) BindReq(c *gin.Context) error {
	return nil
}

func (h *GetFriendRequestListHandler) AfterBindReq() error {
	return nil
}

func (h *GetFriendRequestListHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *GetFriendRequestListHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *GetFriendRequestListHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *GetFriendRequestListHandler) NeedVerifyToken() bool {
	return true
}

func (h *GetFriendRequestListHandler) Process() {
	var err error
	friendRequestList := []*database.FriendRequest{}
	err = database.GetInstance().Model(&database.FriendRequest{}).Where("user_id = ?", h.Req.User.ID).Find(&friendRequestList).Error

	if err != nil {
		msg := fmt.Sprintf("get friend request list error, %v", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		return
	}
	for _, friendRequest := range friendRequestList {
		h.Resp.FriendRequestList = append(h.Resp.FriendRequestList, &detail.RequestDetail{
			ID:       friendRequest.SenderID,
			Name:     friendRequest.Name,
			Avatar:   friendRequest.Avatar,
			Message:  friendRequest.Message,
			Status:   friendRequest.Status,
			CreateAt: friendRequest.CreatedAt,
		})
	}
}
