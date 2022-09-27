package config

import (
	"flag"
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/utils"
	"learn_zinx/Cobra.mayfly/pkg/utils/assert"
	"path/filepath"
)

type Config struct {
	Server *Server `yaml:"server"`
	Jwt    *Jwt    `yaml:"jwt"`
	Aes    *Aes    `yaml:"aes"`
	Redis  *Redis  `yaml:"redis"`
	Mysql  *Mysql  `yaml:"mysql"`
	Log    *Log    `yaml:"log"`
}

// 配置文件内容校验
func (c *Config) Valid() {
	assert.IsTrue(c.Jwt != nil, "配置文件的[jwt]信息不能为空")
	c.Jwt.Valid()
	if c.Aes != nil {
		c.Aes.Valid()
	}
}

type CmdConfigParam struct {
	ConfigFilePath string // -e 配置文件路径
}

// 配置文件映射对象
var Conf *Config

func init() {
	fmt.Println("Init by auto")
}
func Init() {
	configFilePath := flag.String("e", "./config.yml", "配置文件路径，默认可执行文件目录")
	flag.Parse()
	// 从启动参数中，获取配置文件的绝对路径
	path, _ := filepath.Abs(*configFilePath)
	startConfigParam := &CmdConfigParam{ConfigFilePath: path}
	// 读取配置文件信息
	yc := &Config{}
	if err := utils.LoadYml(startConfigParam.ConfigFilePath, yc); err != nil {
		panic(fmt.Sprintf("读取配置文件[%s]失败: %s", startConfigParam.ConfigFilePath, err.Error()))
	}
	// 只要项目验证不通过，就会panic
	yc.Valid()
	Conf = yc
}
