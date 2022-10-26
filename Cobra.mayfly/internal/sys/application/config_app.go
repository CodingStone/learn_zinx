package application

import (
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

type Config interface {
	GetPageList(condition *entity.Config, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Save(config *entity.Config)

	// 获取指定key的配置信息, 不会返回nil, 若不存在则值都默认值即空字符串
	GetConfig(key string) *entity.Config
}

func newConfigApp(configRepo repository.Config) Config {
	return &configAppImpl{
		configRepo: configRepo,
	}
}

type configAppImpl struct {
	configRepo repository.Config
}

func (a *configAppImpl) GetPageList(condition *entity.Config, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return a.configRepo.GetPageList(condition, pageParam, toEntity)
}

func (a *configAppImpl) Save(config *entity.Config) {
	if config.Id == 0 {
		a.configRepo.Insert(config)
	} else {
		a.configRepo.Update(config)
	}
}

func (a *configAppImpl) GetConfig(key string) *entity.Config {
	config := &entity.Config{Key: key}
	if err := a.configRepo.GetConfig(config, "Id", "Key", "Value"); err != nil {
		global.Log.Warnf("不存在key = [%s] 的系统配置", key)
	}
	return config
}
