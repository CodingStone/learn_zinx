package application

import "learn_zinx/Cobra.mayfly/internal/sys/infrastructure/persistence"

var (
	accountApp  = newAccountApp(persistence.GetAccountRepo())
	configApp   = newConfigApp(persistence.GetConfigRepo()) // 项目启动时就会被执行
	msgApp      = newMsgApp(persistence.GetMsgRepo())
	resourceApp = newResourceApp(persistence.GetResourceRepo())
	roleApp     = newRoleApp(persistence.GetRoleRepo())
	syslogApp   = newSyslogApp(persistence.GetSyslogRepo())
)

func GetAccountApp() Account {
	return accountApp
}

// #返回应用查询接口
func GetConfigApp() Config {
	return configApp
}

func GetMsgApp() Msg {
	return msgApp
}

func GetResourceApp() Resource {
	return resourceApp
}

func GetRoleApp() Role {
	return roleApp
}

func GetSyslogApp() Syslog {
	return syslogApp
}
