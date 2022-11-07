package repository

import entity "learn_zinx/Cobra.mayfly/internal/db/domain/entity"

type DbSql interface {
	DeleteBy(condition *entity.DbSql)
}
