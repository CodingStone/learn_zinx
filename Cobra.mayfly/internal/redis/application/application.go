package application

import "learn_zinx/Cobra.mayfly/internal/redis/infrastructure/persistence"

var (
	redisApp Redis = newRedisApp(persistence.GetRedisRepo())
)

func GetRedisApp() Redis {
	return redisApp
}
