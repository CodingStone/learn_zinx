package entity

import "learn_zinx/Cobra.mayfly/pkg/model"

const (
	ConfigKeyUseLoginCaptcha string = "UseLoginCaptcha" // 是否使用登录验证码
)

// # 后面格式为 json 因为要用json形式 进行初始化
type Config struct {
	model.Model        // # 基本数据操作模型
	Name        string `json:"name"` // 配置名
	Key         string `json:"key"`  // 配置key
	Value       string `json:"value"`
	Remark      string `json:"remark"`
}

func (a *Config) TableName() string { // # 表名
	return "t_sys_config"
}

// 若配置信息不存在, 则返回传递的默认值.
// 否则只有value == "1"为true，其他为false
func (c *Config) BoolValue(defaultValue bool) bool {
	// 如果值不存在，则返回默认值
	if c.Id == 0 {
		return defaultValue
	}
	return c.Value == "1"
}
