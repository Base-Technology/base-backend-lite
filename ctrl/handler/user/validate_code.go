package user

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/Base-Technology/base-backend-lite/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

func ValidateCodeHandle(c *gin.Context) {
	hd := &ValidateCodeHandler{}
	handler.Handle(c, hd)
}

type ValidateCodeHandler struct {
	Req  ValidateCodeRequest
	Resp ValidateCodeResponse
}

type ValidateCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type ValidateCodeResponse struct {
	common.BaseResponse
}

func (h *ValidateCodeHandler) BindReq(c *gin.Context) error {
	if err := c.ShouldBindBodyWith(&h.Req, binding.JSON); err != nil {
		msg := fmt.Sprintf("invalid request, bind error: %v", err)
		seelog.Error(msg)
		h.SetError(common.ErrorInvalidParams, msg)
		return err
	}
	return nil
}

func (h *ValidateCodeHandler) AfterBindReq() error {
	return nil
}

func (h *ValidateCodeHandler) GetResponse() interface{} {
	return h.Resp
}

func (h *ValidateCodeHandler) SetError(code int, message string) {
	h.Resp.Code = code
	h.Resp.Message = message
}

func (h *ValidateCodeHandler) SetUser(user *database.User) {
}

func (h *ValidateCodeHandler) NeedVerifyToken() bool {
	return false
}

func (h *ValidateCodeHandler) Process() {
	if !limiter.Take(h.Req.Phone) {
		msg := fmt.Sprintf("rate limit")
		h.SetError(common.ErrorInner, msg)
		return
	}
	code := utils.GenerateValidateCode()
	if err := h.send(code); err != nil {
		msg := fmt.Sprintf("send validate code error, %v", err)
		h.SetError(common.ErrorInner, msg)
		return
	}
	validateCodes.Set(h.Req.Phone, code, cache.DefaultExpiration)
}

func (h *ValidateCodeHandler) send(code string) error {
	request := &SendValidateCodeRequest{
		AppKey:    conf.Conf.ValidateCodeConf.AppKey,
		AppSecret: conf.Conf.ValidateCodeConf.AppSecret,
		AppCode:   conf.Conf.ValidateCodeConf.AppCode,
		Phone:     h.Req.Phone,
		Msg:       fmt.Sprintf("【Base百思】你好，您的验证码为%s", code),
		Timestamp: time.Now().UnixMilli(),
	}

	rawString := fmt.Sprintf("%s%s%d", request.AppKey, request.AppSecret, request.Timestamp)
	request.Sign = fmt.Sprintf("%x", md5.Sum([]byte(rawString)))
	response := &SendValidateCodeResponse{}
	if err := utils.SendHttpRequest(conf.Conf.ValidateCodeConf.Server, http.MethodPost, nil, request, response); err != nil {
		return err
	}
	if response.Code != "00000" {
		return errors.Errorf("%s", response.Desc)
	}
	return nil
}

type SendValidateCodeRequest struct {
	AppKey    string `json:"appkey"`
	AppSecret string `json:"appsecret"`
	AppCode   string `json:"appcode"`
	Phone     string `json:"phone"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
}

type SendValidateCodeResponse struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
}

// TODO: use redis
var validateCodes = cache.New(5*time.Minute, 10*time.Minute)

var limiter = NewLimiter(time.Minute)

type Limiter struct {
	lastTimes map[string]time.Time
	lock      *sync.Mutex
	duration  time.Duration
}

func NewLimiter(duration time.Duration) *Limiter {
	l := &Limiter{
		lastTimes: make(map[string]time.Time),
		lock:      &sync.Mutex{},
		duration:  duration,
	}
	go l.loop()
	return l
}

func (l *Limiter) Take(key string) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	lastTime, ok := l.lastTimes[key]
	if !ok {
		lastTime = time.Time{}
	}
	if !time.Now().After(lastTime.Add(l.duration)) {
		return false
	}
	l.lastTimes[key] = time.Now()
	return true
}

func (l *Limiter) loop() {
	ticker := time.NewTicker(l.duration * 2)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.clear()
			ticker.Reset(l.duration * 2)
		}
	}
}

func (l *Limiter) clear() {
	seelog.Infof("start to clear limiter")
	l.lock.Lock()
	defer l.lock.Unlock()
	clearKeys := []string{}
	for key, value := range l.lastTimes {
		if !time.Now().After(value.Add(l.duration * 2)) {
			clearKeys = append(clearKeys, key)
		}
	}
	for _, key := range clearKeys {
		delete(l.lastTimes, key)
	}
}
