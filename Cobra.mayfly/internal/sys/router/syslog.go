package router

import (
	"learn_zinx/Cobra.mayfly/internal/sys/api"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitSyslogRouter(router *gin.RouterGroup) {
	s := &api.Syslog{
		SyslogApp: application.GetSyslogApp(),
	}
	sys := router.Group("syslogs")
	{
		sys.GET("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(s.Syslogs)
		})
	}
}
