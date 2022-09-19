package zutils

import (
	"../ziface"
	"../zlog"
	"./commandline/args"
	"./commandline/uflag"
	"encoding/json"
	"fmt"
	"os"
)

type GlobalObj struct {
	/*
		Server
	*/
	TCPServer ziface.IServer //当前Zinx的全局Server对象
	Host      string         //当前服务器主机IP
	TCPPort   int            //当前服务器主机监听端口号
	Name      string         //当前服务器名称
	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度

	//Config file path
	ConfFilePath string

	/*
		logger
	*/
	LogDir        string //日志所在文件夹 默认"./log"
	LogFile       string //日志文件名称   默认""  --如果没有设置日志文件，打印信息将打印至stderr
	LogDebugClose bool   //是否关闭Debug日志级别调试信息 默认false  -- 默认打开debug信息
}

// 全局配置对象
var GlobalObject *GlobalObj

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Reload 读取用户的配置文件
func (g *GlobalObj) Reload() {
	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		fmt.Println("Config File ", g.ConfFilePath, " is not exist!!")
		return
	}
	data, err := os.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}

	//将json数据解析到struct中. go中这个方法很方便
	err = json.Unmarshal(data, g)

	if err != nil {
		panic(err)
	}

	//Logger 设置
	if g.LogFile != "" {
		zlog.SetLogFile(g.LogDir, g.LogFile)
	}
	if g.LogDebugClose == true {
		zlog.CloseDebug()
	}
}

// 提供init方法, 默认加载
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		pwd = "."
	}
	// fmt.Println("pwd:", pwd) 程序执行目录 /Cobra/Cobra.Server
	fmt.Println("GlobalObj init~~~", pwd+"/conf/zinx.json")

	// 初始化配置模块flag [其实仅仅是配置，并非执行写入变量操作]
	args.InitConfigFlag(pwd+"/conf/zinx.json", "配置文件，如果没有设置，则默认为<exeDir>/conf/zinx.json")

	// 初始化日志模块flag TODO
	// 解析
	uflag.Parse()
	// 解析之后的操作
	args.FlagHandle()

	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V1.0",
		TCPPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          12000,
		MaxPacketSize:    4096,
		ConfFilePath:     args.Args.ConfigFile,
		WorkerPoolSize:   10,   //10
		MaxWorkerTaskLen: 1024, //1024
		MaxMsgChanLen:    1024, //1024
		LogDir:           pwd + "/log",
		LogFile:          "",
		LogDebugClose:    false,
	}
	//Note: 从配置文件中加载一些用户配置参数
	GlobalObject.Reload()
}
