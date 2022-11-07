package router

import (
	"github.com/gin-gonic/gin"
	"learn_zinx/Cobra.mayfly/internal/common/api"
	dbapp "learn_zinx/Cobra.mayfly/internal/db/application"
	machineapp "learn_zinx/Cobra.mayfly/internal/machine/application"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	redisapp "learn_zinx/Cobra.mayfly/internal/redis/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
)

func InitIndexRouter(router *gin.RouterGroup) {
	index := router.Group("common/index")
	i := &api.Index{
		ProjectApp: projectapp.GetProjectApp(),
		MachineApp: machineapp.GetMachineApp(),
		DbApp:      dbapp.GetDbApp(),
		RedisApp:   redisapp.GetRedisApp(),
	}
	{ // 形成一个作用域，里面会覆盖外边
		// 首页基本信息统计
		index.GET("count", func(g *gin.Context) {
			ctx.NewReqCtxWithGin(g).
				Handle(i.Count)
		})
	}
}
