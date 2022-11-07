package application

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/entity"
	"learn_zinx/Cobra.mayfly/internal/sys/domain/repository"
	"learn_zinx/Cobra.mayfly/pkg/global"
	"learn_zinx/Cobra.mayfly/pkg/model"
)

// # 定义查询接口
type Config interface {
	GetPageList(condition *entity.Config, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Save(config *entity.Config)

	// 获取指定key的配置信息, 不会返回nil, 若不存在则值都默认值即空字符串
	GetConfig(key string) *entity.Config
}

type configAppImpl struct {
	configRepo repository.Config
}

// # 对Conf接口实现（也可以理解为对基础查询方法封装，实现贴合业务查询）
func newConfigApp(configRepo repository.Config) Config {
	return &configAppImpl{
		configRepo: configRepo,
	}
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
	fmt.Printf("执行到这里: %v \n", key)
	config := &entity.Config{Key: key}
	fmt.Printf("执行到这里: %+v \n", config)

	// # 这里比较绕。 因为跳转地方是 interface 类型。 具体实现文件路径为 internal/sys/infrastructure/persistence/config_repo.go
	if err := a.configRepo.GetConfig(config, "Id", "Key", "Value"); err != nil { // # 要到DB中取数据
		fmt.Printf("取数据错误为: %+v\n", err)
		global.Log.Warnf("不存在key = [%s] 的系统配置", key)
	}
	return config
}
