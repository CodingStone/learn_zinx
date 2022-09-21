package zrouter

import (
	"fmt"
	"learn_zinx/Cobra.Server/ziface"
	"learn_zinx/Cobra.Server/znet"
)

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//Ping Handle
func (this *PingRouter) Handle(request ziface.IRequest) {

	fmt.Println("Call PingRouter Handle")
	//zlog.Debug("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))
	//zlog.Debug("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendBuffMsg(0, []byte("ping...ping...ping[FromServer]"))
	if err != nil {
		fmt.Println(err)
		//zlog.Error(err)
	}
}