package repository

import (
	entity "learn_zinx/Cobra.mayfly/internal/db/domain/entry"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type Db interface {
	// 分页获取机器信息列表
	GetDbList(condition *entity.Db, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Count(condition *entity.Db) int64

	// 根据条件获取账号信息
	GetDb(condition *entity.Db, cols ...string) error

	// 根据id获取
	GetById(id uint64, cols ...string) *entity.Db

	Insert(db *entity.Db)

	Update(db *entity.Db)

	Delete(id uint64)
}
