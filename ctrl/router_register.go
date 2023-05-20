package ctrl

import (
	"github.com/Base-Technology/base-app-lite/common"
	"github.com/Base-Technology/base-app-lite/ctrl/handler/post"
	"github.com/Base-Technology/base-app-lite/ctrl/handler/user"
)

func init() {
	initUserInterfaces()
	initPostInterfaces()
}

func initUserInterfaces() {
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/register", Method: "POST", Handler: user.RegisterHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/login", Method: "POST", Handler: user.LoginHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/reset_password", Method: "POST", Handler: user.ResetPasswordHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/info", Method: "GET", Handler: user.GetInfoHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/info", Method: "PUT", Handler: user.SetInfoHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/users", Method: "GET", Handler: user.GetUserHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/follow", Method: "POST", Handler: user.FollowHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/follow", Method: "DELETE", Handler: user.CancelFollowHandle})
}

func initPostInterfaces() {
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts", Method: "POST", Handler: post.CreatePostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts", Method: "DELETE", Handler: post.DeletePostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts", Method: "GET", Handler: post.GetPostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts/like", Method: "POST", Handler: post.LikePostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts/like", Method: "DELETE", Handler: post.UnlikePostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts/collect", Method: "POST", Handler: post.CollectPostHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/posts/collect", Method: "DELETE", Handler: post.UncollectPostHandle})
}
