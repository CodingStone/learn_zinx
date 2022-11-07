package persistence

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/biz"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type configRepoImpl struct{}

// # 对配置表的接口实现 为什么实现 会写在这里
func newConfigRepo() repository.Config {
	return new(configRepoImpl)
}

func (m *configRepoImpl) GetPageList(condition *entity.Config, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return model.GetPage(pageParam, condition, toEntity)
}

func (m *configRepoImpl) Insert(config *entity.Config) {
	biz.ErrIsNil(model.Insert(config), "新增系统配置失败")
}

func (m *configRepoImpl) Update(config *entity.Config) {
	biz.ErrIsNil(model.UpdateById(config), "更新系统配置失败")
}

// # 没太理解 为什么实现要放在这里
func (m *configRepoImpl) GetConfig(condition *entity.Config, cols ...string) error {
	fmt.Printf("配置条件为: %+v\n", condition)
	return model.GetBy(condition, cols...)
}

func (r *configRepoImpl) GetByCondition(condition *entity.Config, cols ...string) error {
	return model.GetBy(condition, cols...)
}
