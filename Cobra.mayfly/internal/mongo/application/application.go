package application

import "learn_zinx/Cobra.mayfly/internal/mongo/infrastructure/persistence"

var (
	mongoApp Mongo = newMongoAppImpl(persistence.GetMongoRepo())
)

func GetMongoApp() Mongo {
	return mongoApp
}
