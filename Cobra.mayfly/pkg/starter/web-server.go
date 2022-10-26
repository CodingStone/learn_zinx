package starter

import (
	"learn_zinx/Cobra.mayfly/initialize"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
)

func runWebServer() {
	println("run web server~~~")
	// 权限处理器 [注册相应函数]
	ctx.UseBeforeHandlerInterceptor(ctx.PermissionHandler)
	// 日志处理器
	ctx.UseAfterHandlerInterceptor(ctx.LogHandler)
	// 设置日志保存函数
	ctx.SetSaveLogFunc(initialize.InitSaveLogFunc())
}
