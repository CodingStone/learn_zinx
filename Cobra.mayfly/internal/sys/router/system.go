package router

import (
	"learn_zinx/Cobra.mayfly/internal/sys/api"

	"github.com/gin-gonic/gin"
)

func InitSystemRouter(router *gin.RouterGroup) {
	s := &api.System{}
	sys := router.Group("sysmsg")

	{
		sys.GET("", s.ConnectWs)
	}
}
