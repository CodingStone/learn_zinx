package main

import (
	"../../ziface"
	"../../znet"
	"./zrouter"
)


//创建连接的时候执行
func DoConnectionBegin(conn ziface.IConnection) {
	println("DoConnectionBegin ~~~")
	//zlog.Debug("DoConnecionBegin is Called ... ")
	//
	////设置两个链接属性，在连接创建之后
	//zlog.Debug("Set conn Name, Home done!")
	//conn.SetProperty("Name", "Aceld")
	//conn.SetProperty("Home", "https://www.kancloud.cn/@aceld")
	//
	//err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	//if err != nil {
	//	zlog.Error(err)
	//}
}

//连接断开的时候执行
func DoConnectionLost(conn ziface.IConnection) {
	println("DoConnectionLost ~~~")
	//在连接销毁之前，查询conn的Name，Home属性
	//if name, err := conn.GetProperty("Name"); err == nil {
	//	zlog.Error("Conn Property Name = ", name)
	//}
	//
	//if home, err := conn.GetProperty("Home"); err == nil {
	//	zlog.Error("Conn Property Home = ", home)
	//}
	//
	//zlog.Debug("DoConneciotnLost is Called ... ")
}

func main()  {
	s := znet.NewServer()
	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &zrouter.PingRouter{})
	s.AddRouter(1, &zrouter.HelloZinxRouter{})

	//开启服务
	s.Serve()

}
