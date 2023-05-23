package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/database"
	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/Base-Technology/base-backend-lite/token"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	BindReq(c *gin.Context) error
	AfterBindReq() error
	GetResponse() interface{}
	SetError(code int, message string)
	SetUser(user *database.User)
	NeedVerifyToken() bool
	Process()
}

func Handle(c *gin.Context, hd Handler) {
	defer func() {
		if err := recover(); err != nil {
			hd.SetError(common.ErrorPanic, fmt.Sprintf("%v", err))
			seelog.Errorf("panic: %v", err)
			stack := debug.Stack()
			seelog.Errorf("panic stack: %s", string(stack))
		}
		c.JSON(http.StatusOK, hd.GetResponse())
	}()
	if hd.NeedVerifyToken() {
		if err := verifyToken(c, hd); err != nil {
			return
		}
	}
	if err := hd.BindReq(c); err != nil {
		return
	}
	if err := hd.AfterBindReq(); err != nil {
		return
	}
	hd.Process()
}

func verifyToken(c *gin.Context, hd Handler) error {
	auth := c.Request.Header.Get("Authorization")
	kv := strings.Split(auth, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		msg := "invalid Authorization"
		seelog.Error(msg)
		hd.SetError(common.ErrorInvalidToken, msg)
		return fmt.Errorf(msg)
	}
	id, err := token.VerifyToken(kv[1])
	if err != nil {
		msg := fmt.Sprintf("verify token error, %v", err)
		seelog.Error(msg)
		hd.SetError(common.ErrorInvalidToken, msg)
		return err
	}
	user := &database.User{}
	if err := database.GetInstance().First(user, id).Error; err != nil {
		msg := fmt.Sprintf("get user error, %v", err)
		hd.SetError(common.ErrorInvalidToken, msg)
		return err
	}
	hd.SetUser(user)
	return nil
}
