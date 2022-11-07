package repository

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

// # 基本操作约束、相当于定义接口
type Config interface {
	GetPageList(condition *entity.Config, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Insert(config *entity.Config)

	Update(config *entity.Config)

	GetConfig(config *entity.Config, cols ...string) error

	GetByCondition(condition *entity.Config, cols ...string) error
}
