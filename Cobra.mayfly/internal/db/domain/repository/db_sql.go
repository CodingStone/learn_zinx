package repository

import entity "learn_zinx/Cobra.mayfly/internal/db/domain/entry"

type DbSql interface {
	DeleteBy(condition *entity.DbSql)
}
