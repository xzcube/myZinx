package znet

import (
	"myZinx/ziface"
)

type Request struct {
	// 已经和客户端建立好的连接
	conn ziface.IConnections

	// 客户端请求的数据
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnections {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func(r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}


