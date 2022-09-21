package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"learn_zinx/Cobra.Server/znet"
)

/*
	模拟客户端
*/
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return
	}

	for {
		//发封包message消息
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMsgPackage(0, []byte("Zinx client Demo Test MsgID=0, [Ping]")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull 会把msg填充满为止【相当于把消息ID和数据长度】
		if err != nil {
			fmt.Println("read head error", err)
			break
		}
		//将headData字节流 拆包到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg 是有data数据的，需要再次读取data数据
			// s, ok := x.(T)
			//如果 x 不是 nil，且 x 可以转换成 T 类型，就会断言成功，返回 T 类型的变量 s。
			//如果 T 不是接口类型，则要求 x 的类型就是 T，如果 T 是一个接口，要求 x 实现了 T 接口。
			//如果断言类型成立，则表达式返回值就是 T 类型的 x，如果断言失败就会触发 panic。


			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data) // 会接着从conn中读数据，写入到msg.Data中
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}

			fmt.Println("==> Test Router:[Ping] Recv Msg: ID=", msg.ID, ", len=", msg.DataLen, ", data=", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}
}
