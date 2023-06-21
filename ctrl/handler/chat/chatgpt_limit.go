package chat

import (
	//"encoding/json"
	//"fmt"

	"time"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/ctrl/types"
	"github.com/Base-Technology/base-backend-lite/database"

	//"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
)

func ChatGPTLimitHandle(c *gin.Context) {
	hd := &ChatGPTLimitHandler{}
	handler.Handle(c, hd)
}

type ChatGPTLimitHandler struct {
	Req  ChatGPTLimitRequest
	Resp ChatGPTLimitResponse
}

type ChatGPTLimitRequest struct {
	User *database.User
}

type ChatGPTLimitResponse struct {
	common.BaseResponse
	types.ChatGPTLimitDetail
}

func (h *ChatGPTLimitHandler) BindReq(c *gin.Context) error {
	return nil
}

func (h *ChatGPTLimitHandler) AfterBindReq() error {
	return nil
}

func (h *ChatGPTLimitHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *ChatGPTLimitHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *ChatGPTLimitHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *ChatGPTLimitHandler) NeedVerifyToken() bool {
	return true
}

func (h *ChatGPTLimitHandler) Process() {
	limit := &database.ChatGPTLimit{}
	database.GetInstance().
		Where(database.ChatGPTLimit{UserID: h.Req.User.ID}).
		Attrs(database.ChatGPTLimit{LastResetTime: time.Now()}).
		FirstOrCreate(limit)

	// reset limit if last reset time is more than 24 hours
	if time.Since(limit.LastResetTime).Hours() > 24 {
		limit.DailyLeftCallCount = limit.MaxDailyCallCount
		limit.DailyLeftTokenCount = limit.MaxDailyTokenCount
		limit.LastResetTime = time.Now()
	}

	h.Resp.ChatGPTLimitDetail = types.ChatGPTLimitDetail{
		DailyLeftCallCount:  limit.DailyLeftCallCount,
		DailyLeftTokenCount: limit.DailyLeftTokenCount,
		TotalTokenLeftCount: limit.TotalTokenLeftCount,
		MaxDailyCallCount:   limit.MaxDailyCallCount,
		MaxDailyTokenCount:  limit.MaxDailyTokenCount,
	}
}
