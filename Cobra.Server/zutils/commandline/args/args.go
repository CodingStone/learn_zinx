package args

import (
	"../uflag"
	"os"
	"path"
)

type args struct {
	ExeAbsDir  string
	ExeName    string
	ConfigFile string
}

var (
	Args   = args{}
	isInit = false
)

// 初始化获得程序执行名称、与程序执行路径
func init() {
	exe := os.Args[0]

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	Args.ExeAbsDir = pwd
	Args.ExeName = path.Base(exe)
}

// 根据配置文件初始化项目。defaultValue: 配置文件路径。 tips: 文件找不到提示
func InitConfigFlag(defaultValue string, tips string) {
	if isInit {
		return
	}
	isInit = true
	//如果 程序通过 -c 参数配置，那么写入到地址 Args.ConfigFile 中，defalutValues 是默认值, tips是提示
	uflag.StringVar(&Args.ConfigFile, "c", defaultValue, tips) // 这块仅仅是配置，需要调用parse()方法进行写入操作
	return
}

// 判断 是不是 相对地址，补全地址
func FlagHandle() {
	if !path.IsAbs(Args.ConfigFile) {
		Args.ConfigFile = path.Join(Args.ExeAbsDir, Args.ConfigFile)
	}
}
