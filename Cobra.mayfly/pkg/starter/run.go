package starter

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/config"
	"learn_zinx/Cobra.mayfly/pkg/logger"
)

func RunWebServer() {
	fmt.Println("Run Web Server~~~")
	// 初始化项目配置，从yaml中读取配置项
	config.Init()

	// 初始化日志配置信息
	logger.Init()
}
