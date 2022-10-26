package repository

import (
	"learn_zinx/Cobra.mayfly/internal/mongo/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type Mongo interface {
	// 分页获取列表
	GetList(condition *entity.Mongo, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Count(condition *entity.Mongo) int64

	// 根据条件获取
	Get(condition *entity.Mongo, cols ...string) error

	// 根据id获取
	GetById(id uint64, cols ...string) *entity.Mongo

	Insert(db *entity.Mongo)

	Update(db *entity.Mongo)

	Delete(id uint64)
}
