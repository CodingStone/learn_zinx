package starter

import (
	"learn_zinx/Cobra.mayfly/initialize"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
	"learn_zinx/Cobra.mayfly/pkg/global"
)

func runWebServer() {
	println("run web server~~~")
	// 权限处理器 [注册相应函数]
	ctx.UseBeforeHandlerInterceptor(ctx.PermissionHandler)
	// 日志处理器
	ctx.UseAfterHandlerInterceptor(ctx.LogHandler)
	// 设置日志保存函数
	ctx.SetSaveLogFunc(initialize.InitSaveLogFunc())

	// 注册路由
	web := initialize.InitRouter()
	server := config.Conf.Server
	port := server.GetPort()
	global.Log.Infof("Listening and serving HTTP on %s", port)

	var err error
	if server.Tls != nil && server.Tls.Enable {
		err = web.RunTLS(port, server.Tls.CertFile, server.Tls.KeyFile)
	} else {
		err = web.Run(port)
	}
	biz.ErrIsNilAppendErr(err, "服务启动失败: %s")
}
