package router

import (
	"learn_zinx/Cobra.mayfly/internal/sys/api"
	"learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitRoleRouter(router *gin.RouterGroup) {
	r := &api.Role{
		RoleApp:     application.GetRoleApp(),
		ResourceApp: application.GetResourceApp(),
	}
	db := router.Group("sys/roles")
	{ // 仅仅是相当于代码块作用

		db.GET("", func(c *gin.Context) { // 获得路径
			ctx.NewReqCtxWithGin(c).Handle(r.Roles)
		})

		saveRole := ctx.NewLogInfo("保存角色").WithSave(true)
		sPermission := ctx.NewPermission("role:add")
		db.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(saveRole).
				WithRequiredPermission(sPermission).
				Handle(r.SaveRole)
		})

		delRole := ctx.NewLogInfo("删除角色").WithSave(true)
		drPermission := ctx.NewPermission("role:del")
		db.DELETE(":id", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(delRole).
				WithRequiredPermission(drPermission).
				Handle(r.DelRole)
		})

		db.GET(":id/resourceIds", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.RoleResourceIds)
		})

		db.GET(":id/resources", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(r.RoleResource)
		})

		saveResource := ctx.NewLogInfo("保存角色资源").WithSave(true)
		srPermission := ctx.NewPermission("role:saveResources")
		db.POST(":id/resources", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(saveResource).
				WithRequiredPermission(srPermission).
				Handle(r.SaveResource)
		})
	}
}
