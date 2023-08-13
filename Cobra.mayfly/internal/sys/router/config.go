package router

import (
	"learn_zinx/Cobra.mayfly/internal/sys/api"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitSysConfigRouter(router *gin.RouterGroup) {
	//# r：表示API请求，里面包含具体请求对应处理方法 【后边参数是数据库操作方法集合】
	r := &api.Config{ConfigApp: application.GetConfigApp()}
	db := router.Group("sys/configs")
	{
		db.GET("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.Configs)
		})

		// # 是否开启登陆验证
		db.GET("/value", func(c *gin.Context) {
			// # 请求设计为链式操作。创建上下文，设置token，并将 路由 要处理的方法传入进来
			ctx.NewReqCtxWithGin(c).WithNeedToken(false).Handle(r.GetConfigValueByKey)
		})

		saveConfig := ctx.NewLogInfo("保存系统配置信息").WithSave(true)
		db.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).
				WithLog(saveConfig).
				Handle(r.SaveConfig)
		})
	}
}
