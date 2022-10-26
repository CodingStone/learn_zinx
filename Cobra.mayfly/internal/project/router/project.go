package router

import (
	"learn_zinx/Cobra.mayfly/internal/project/api"
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
	sysapp "learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"

	"github.com/gin-gonic/gin"
)

func InitProjectRouter(router *gin.RouterGroup) {
	m := &api.Project{
		ProjectApp: projectapp.GetProjectApp(),
		AccountApp: sysapp.GetAccountApp()}

	project := router.Group("/projects")
	{
		router.GET("/accounts/projects", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(m.GetProjectsByLoginAccount)
		})

		// 获取项目列表
		project.GET("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(m.GetProjects)
		})

		saveProjectLog := ctx.NewLogInfo("保存项目信息").WithSave(true)
		savePP := ctx.NewPermission("project:save")
		// 保存项目下的环境信息
		project.POST("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(saveProjectLog).
				WithRequiredPermission(savePP).
				Handle(m.SaveProject)
		})

		delProjectLog := ctx.NewLogInfo("删除项目信息").WithSave(true)
		delPP := ctx.NewPermission("project:del")
		// 删除项目
		project.DELETE("", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(delProjectLog).
				WithRequiredPermission(delPP).
				Handle(m.DelProject)
		})

		// 获取项目下的环境信息列表
		project.GET("/:projectId/envs", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(m.GetProjectEnvs)
		})

		saveProjectEnvLog := ctx.NewLogInfo("新增项目环境信息").WithSave(true)
		savePeP := ctx.NewPermission("project:env:add")
		// 保存项目下的环境信息
		project.POST("/:projectId/envs", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(saveProjectEnvLog).
				WithRequiredPermission(savePeP).
				Handle(m.SaveProjectEnvs)
		})

		delProjectEnvLog := ctx.NewLogInfo("删除项目环境信息").WithSave(true)
		project.DELETE("/envs", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(delProjectEnvLog).
				WithRequiredPermission(delPP).
				Handle(m.DelProjectEnv)
		})

		// 获取项目下的成员信息列表
		project.GET("/:projectId/members", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).Handle(m.GetProjectMembers)
		})

		// 保存项目成员
		saveProjectMemLog := ctx.NewLogInfo("新增项目成员").WithSave(true)
		savePmP := ctx.NewPermission("project:member:add")
		project.POST("/:projectId/members", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(saveProjectMemLog).
				WithRequiredPermission(savePmP).
				Handle(m.SaveProjectMember)
		})

		// 删除项目成员
		delProjectMemLog := ctx.NewLogInfo("删除项目成员").WithSave(true)
		savePmdP := ctx.NewPermission("project:member:del")
		project.DELETE("/:projectId/members/:accountId", func(c *gin.Context) {
			ctx.NewReqCtxWithGin(c).WithLog(delProjectMemLog).
				WithRequiredPermission(savePmdP).
				Handle(m.DelProjectMember)
		})
	}
}
