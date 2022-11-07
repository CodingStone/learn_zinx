package router

import "github.com/gin-gonic/gin"

func Init(router *gin.RouterGroup) {
	InitCaptchaRouter(router) // # 注册验证码路由
	InitAccountRouter(router) // 注册account路由 # 里面包含登陆路由
	InitResourceRouter(router)
	InitRoleRouter(router)
	InitSystemRouter(router)
	InitSyslogRouter(router)
	InitSysConfigRouter(router)
}
