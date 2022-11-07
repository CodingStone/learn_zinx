package persistence

import "learn_zinx/Cobra.mayfly/internal/sys/domain/repository"

var (
	accountRepo  = newAccountRepo()
	configRepo   = newConfigRepo() // 配置表 t_sys_config 查询实现
	msgRepo      = newMsgRepo()
	resourceRepo = newResourceRepo()
	roleRepo     = newRoleRepo()
	syslogRepo   = newSyslogRepo()
)

func GetAccountRepo() repository.Account {
	return accountRepo
}

func GetConfigRepo() repository.Config {
	return configRepo
}

func GetMsgRepo() repository.Msg {
	return msgRepo
}

func GetResourceRepo() repository.Resource {
	return resourceRepo
}

func GetRoleRepo() repository.Role {
	return roleRepo
}

func GetSyslogRepo() repository.Syslog {
	return syslogRepo
}
