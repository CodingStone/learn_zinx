package persistence

import "learn_zinx/Cobra.mayfly/internal/db/domain/repository"

var (
	dbRepo        repository.Db        = newDbRepo()
	dbSqlRepo     repository.DbSql     = newDbSqlRepo()
	dbSqlExecRepo repository.DbSqlExec = newDbSqlExecRepo()
)

func GetDbRepo() repository.Db {
	return dbRepo
}

func GetDbSqlRepo() repository.DbSql {
	return dbSqlRepo
}

func GetDbSqlExecRepo() repository.DbSqlExec {
	return dbSqlExecRepo
}
