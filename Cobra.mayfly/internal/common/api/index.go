package api

import (
	projectapp "learn_zinx/Cobra.mayfly/internal/project/application"
)

type Index struct {
	ProjectApp projectapp.Project
	MachineApp machineapp.Machine
	DbApp      dbapp.Db
	RedisApp   redisapp.Redis
}
