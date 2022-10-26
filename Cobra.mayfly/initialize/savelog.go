package initialize

import (
	sysapp "learn_zinx/Cobra.mayfly/internal/sys/application"
	"learn_zinx/Cobra.mayfly/pkg/ctx"
)

func InitSaveLogFunc() ctx.SaveLogFunc {
	return sysapp.GetSyslogApp().SaveFromReq
}
