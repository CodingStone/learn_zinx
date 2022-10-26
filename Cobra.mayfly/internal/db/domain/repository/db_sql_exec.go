package repository

import (
	"learn_zinx/Cobra.mayfly/internal/db/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type DbSqlExec interface {
	Insert(d *entity.DbSqlExec)

	DeleteBy(condition *entity.DbSqlExec)

	// 分页获取
	GetPageList(condition *entity.DbSqlExec, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult
}
