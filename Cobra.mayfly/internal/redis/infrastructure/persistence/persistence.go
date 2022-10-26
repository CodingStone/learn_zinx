package persistence

import "learn_zinx/Cobra.mayfly/internal/redis/domain/repository"

var (
	redisRepo repository.Redis = newRedisRepo()
)

func GetRedisRepo() repository.Redis {
	return redisRepo
}
