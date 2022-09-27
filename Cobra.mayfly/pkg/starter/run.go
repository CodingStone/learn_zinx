package starter

import (
	"fmt"
	"learn_zinx/Cobra.mayfly/pkg/config"
)

func RunWebServer() {
	fmt.Println("Run Web Server~~~")
	// 初始化项目配置，从yaml中读取配置项
	config.Init()
}
