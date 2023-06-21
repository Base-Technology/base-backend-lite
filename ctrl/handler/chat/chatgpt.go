package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/ctrl/types"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func ChatGPTHandle(c *gin.Context) {
	hd := &ChatGPTHandler{}
	handler.Handle(c, hd)
}

type ChatGPTHandler struct {
	Req  ChatGPTRequest
	Resp ChatGPTResponse
}

type ChatGPTRequest struct {
	Prompt string `json:"prompt" binding:"required"`
	User   *database.User
}

type ChatGPTResponse struct {
	common.BaseResponse
	Response string `json:"response"`
	types.ChatGPTLimitDetail
}

type ChatGPTProxyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func (h *ChatGPTHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *ChatGPTHandler) AfterBindReq() error {
	return nil
}

func (h *ChatGPTHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *ChatGPTHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *ChatGPTHandler) SetUser(user *database.User) {
	h.Req.User = user
}

func (h *ChatGPTHandler) NeedVerifyToken() bool {
	return true
}

func (h *ChatGPTHandler) Process() {
	limit := &database.ChatGPTLimit{}
	database.GetInstance().
		Where(database.ChatGPTLimit{UserID: h.Req.User.ID}).
		Attrs(database.ChatGPTLimit{LastResetTime: time.Now()}).
		FirstOrCreate(limit)

	// reset limit if last reset time is more than 24 hours
	resetBalance(limit)

	// fill in limit first, if there is error, this will be the response
	h.Resp.ChatGPTLimitDetail = types.ChatGPTLimitDetail{
		DailyLeftCallCount:  limit.DailyLeftCallCount,
		DailyLeftTokenCount: limit.DailyLeftTokenCount,
		TotalTokenLeftCount: limit.TotalTokenLeftCount,
		MaxDailyCallCount:   limit.MaxDailyCallCount,
		MaxDailyTokenCount:  limit.MaxDailyTokenCount,
	}

	if !enoughBalance(limit, &h.Req.Prompt) {
		seelog.Errorf("limit exceedeed: %+v", limit)
		h.SetError(common.ErrorLimitExceedeed, "limit exceedeed")
		return
	}

	url := fmt.Sprintf("%s?prompt=%s",
		conf.Conf.ChatGPTProxyConf.Url,
		url.QueryEscape(h.Req.Prompt))
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		seelog.Errorf("http get error: %v", err)
		h.SetError(common.ErrorInner, "http get error")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Errorf("read response body error: %v", err)
		h.SetError(common.ErrorInner, "read response body error")
		return
	}

	proxy_resp := &ChatGPTProxyResponse{}
	if err = json.Unmarshal(body, proxy_resp); err != nil {
		seelog.Errorf("unmarshal response body error: %v, %v", err, string(body))
		h.SetError(common.ErrorInner, "unmarshal response body error")
		return
	}
	//fmt.Println(proxy_resp.Data)
	h.Resp.Response = proxy_resp.Data

	limit = &database.ChatGPTLimit{}
	tx := database.GetInstance().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err := tx.Commit().Error; err != nil {
			msg := fmt.Sprintf("ChatGPT: UserId=%v commit error [%v]", h.Req.User.ID, err)
			seelog.Errorf(msg)
			h.SetError(common.ErrorInner, msg)
			tx.Rollback()
		}
	}()
	if err := tx.Where(database.ChatGPTLimit{UserID: h.Req.User.ID}).First(limit).Error; err != nil {
		msg := fmt.Sprintf("ChatGPT: UserId=%v not found [%v]", h.Req.User.ID, err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}
	resetBalance(limit)
	updateBalance(limit, &proxy_resp.Data)
	if err := tx.Save(limit).Error; err != nil {
		msg := fmt.Sprintf("ChatGPT: error when saving limit [%v]", err)
		seelog.Errorf(msg)
		h.SetError(common.ErrorInner, msg)
		tx.Rollback()
		return
	}

	// fill in the updated limit
	h.Resp.ChatGPTLimitDetail = types.ChatGPTLimitDetail{
		DailyLeftCallCount:  limit.DailyLeftCallCount,
		DailyLeftTokenCount: limit.DailyLeftTokenCount,
		TotalTokenLeftCount: limit.TotalTokenLeftCount,
		MaxDailyCallCount:   limit.MaxDailyCallCount,
		MaxDailyTokenCount:  limit.MaxDailyTokenCount,
	}
}
