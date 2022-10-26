package application

import (
	dbapp "learn_zinx/Cobra.mayfly/internal/db/application"
	machineapp "learn_zinx/Cobra.mayfly/internal/machine/application"
	mongoapp "learn_zinx/Cobra.mayfly/internal/mongo/application"
	"learn_zinx/Cobra.mayfly/internal/project/infrastructure/persistence"
	redisapp "learn_zinx/Cobra.mayfly/internal/redis/application"
)

var (
	projectApp Project = newProjectApp(
		persistence.GetProjectRepo(),
		persistence.GetProjectEnvRepo(),
		persistence.GetProjectMemberRepo(),
		machineapp.GetMachineApp(),
		redisapp.GetRedisApp(),
		dbapp.GetDbApp(),
		mongoapp.GetMongoApp(),
	)
)

func GetProjectApp() Project {
	return projectApp
}
