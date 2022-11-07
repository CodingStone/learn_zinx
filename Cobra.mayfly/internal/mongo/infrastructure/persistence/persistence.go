package persistence

import (
	"learn_zinx/Cobra.mayfly/internal/mongo/domain/repository"
)

var (
	mongoRepo repository.Mongo = newMongoRepo()
)

func GetMongoRepo() repository.Mongo {
	return mongoRepo
}
