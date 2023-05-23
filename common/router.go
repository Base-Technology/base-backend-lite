package common

import (
	"strings"

	"github.com/Base-Technology/base-backend-lite/seelog"
	"github.com/gin-gonic/gin"
)

var RouterRegister RouterHandlerRegister

type RouterHandlerRegister struct {
	Router         *gin.Engine
	RouterHandlers []RouterHandler
}

type RouterHandler struct {
	Path    string          `json:"path"`
	Method  string          `json:"method"`
	Handler gin.HandlerFunc `json:"handler"`
}

func (r *RouterHandlerRegister) SetRouter(router *gin.Engine) {
	r.Router = router
}

func (r *RouterHandlerRegister) InitRouter() {
	seelog.Infof("init router")
	for _, routerHandler := range r.RouterHandlers {
		switch strings.ToUpper(routerHandler.Method) {
		case "GET":
			r.Router.GET(routerHandler.Path, routerHandler.Handler)
			seelog.Infof("register GET router: %s", routerHandler.Path)
		case "POST":
			r.Router.POST(routerHandler.Path, routerHandler.Handler)
			seelog.Infof("register POST router: %s", routerHandler.Path)
		case "PUT":
			r.Router.PUT(routerHandler.Path, routerHandler.Handler)
			seelog.Infof("register PUT router: %s", routerHandler.Path)
		case "DELETE":
			r.Router.DELETE(routerHandler.Path, routerHandler.Handler)
			seelog.Infof("register DELETE router: %s", routerHandler.Path)
		default:
			seelog.Errorf("not support method %s for path %s", routerHandler.Method, routerHandler.Path)
		}
	}
}

func (r *RouterHandlerRegister) RegisterRouterHandler(rh RouterHandler) {
	if r.RouterHandlers == nil {
		r.RouterHandlers = make([]RouterHandler, 0)
	}
	r.RouterHandlers = append(r.RouterHandlers, rh)
}
