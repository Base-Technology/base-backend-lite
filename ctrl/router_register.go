package ctrl

import (
	"github.com/Base-Technology/base-backend-lite/common"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler/chat"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler/group"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler/post"
	"github.com/Base-Technology/base-backend-lite/ctrl/handler/user"
)

func init() {
	initUserInterfaces()
	initPostInterfaces()
	initGroupInterfaces()
	initChatInterfaces()
}

func initUserInterfaces() {
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/register", Method: "POST", Handler: user.RegisterHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/login", Method: "POST", Handler: user.LoginHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/reset_password", Method: "POST", Handler: user.ResetPasswordHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/info", Method: "GET", Handler: user.GetInfoHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/info", Method: "PUT", Handler: user.SetInfoHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/user", Method: "GET", Handler: user.GetOtherUserHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/users", Method: "GET", Handler: user.GetUserHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/follow", Method: "POST", Handler: user.FollowHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/follow", Method: "DELETE", Handler: user.CancelFollowHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/validate_code", Method: "POST", Handler: user.ValidateCodeHandle})
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

func initGroupInterfaces() {
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/group/user", Method: "GET", Handler: group.GetGroupUserHandle})
}

func initChatInterfaces() {
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/chat/chatgpt", Method: "POST", Handler: chat.ChatGPTHandle})
	common.RouterRegister.RegisterRouterHandler(common.RouterHandler{Path: "/api/v1/chat/chatgpt_limit", Method: "GET", Handler: chat.ChatGPTLimitHandle})
}
