package router

import (
	"github.com/gin-gonic/gin"
	"learn_zinx/Cobra.mayfly/internal/common/api"
)

func InitIndexRouter(router *gin.RouterGroup) {
	index := router.Group("common/index")
	i := &api.Index{
		ProjectApp: projectapp.GetProjectApp(),
		MachineApp: machineapp.GetMachineApp(),
		DbApp:      dbapp.GetDbApp(),
		RedisApp:   redisapp.GetRedisApp(),
	}
}
