package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const ChatGPTProxyURL = "http://147.182.251.92:5000/proxy/openai"

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
	ChatGPTLimit
}

type ChatGPTLimit struct {
	UsedCallCount       int `json:"used_call_count"`
	LeftCallCount       int `json:"left_call_count"`
	DailyUsedTokenCount int `json:"daily_used_token_count"`
	DailyLeftTokenCount int `json:"daily_left_token_count"`
	TotalTokenCount     int `json:"total_token_count"`
	TotalTokenLeftCount int `json:"total_token_left_count"`
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
	// TODO: enable authorization
	return false
}

func (h *ChatGPTHandler) Process() {
	url := fmt.Sprintf("%s?prompt=%s", ChatGPTProxyURL, url.QueryEscape(h.Req.Prompt))
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		seelog.Errorf("http get error: %v", err)
		h.SetError(common.ErrorInner, "internal error")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		seelog.Errorf("read response body error: %v", err)
		h.SetError(common.ErrorInner, "internal error")
		return
	}

	proxy_resp := &ChatGPTProxyResponse{}
	if err := json.Unmarshal(body, proxy_resp); err != nil {
		seelog.Errorf("unmarshal response body error: %v", err)
		h.SetError(common.ErrorInner, "internal error")
		return
	}
	//fmt.Println(proxy_resp.Data)
	h.Resp.Response = proxy_resp.Data
}
