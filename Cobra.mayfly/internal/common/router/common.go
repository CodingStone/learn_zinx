package router

import (
	"learn_zinx/Cobra.mayfly/internal/common/api"
	"learn_zinx/Cobra.mayfly/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitCommonRouter(router *gin.RouterGroup) {
	common := router.Group("common")
	c := &api.Common{}
	{
		// # 获取公钥 加密密码用
		common.GET("public-key", func(g *gin.Context) {
			ctx.NewReqCtxWithGin(g).
				WithNeedToken(false).
				Handle(c.RasPublicKey)
		})
	}
}
