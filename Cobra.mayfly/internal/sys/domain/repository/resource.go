package repository

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
)

type Resource interface {
	// 获取资源列表
	GetResourceList(condition *entity.Resource, toEntity interface{}, orderBy ...string)

	GetById(id uint64, cols ...string) *entity.Resource

	GetByIdIn(ids []uint64, toEntity interface{}, orderBy ...string)

	Delete(id uint64)

	GetByCondition(condition *entity.Resource, cols ...string) error

	// 获取账号资源列表
	GetAccountResources(accountId uint64, toEntity interface{})
}
