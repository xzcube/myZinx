package znet

import (
	"fmt"
	"myZinx/utils"
	"myZinx/ziface"
	"net"
)

// iServer的接口实现，定义一个Server服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	// 服务器的监听ip
	IP string
	// 服务器的监听端口
	Port int
	// 当前server的消息管理模块，用来绑定msgID和对应的业务处理api
	MsgHandle ziface.IMsgHandle

	// 该Server的连接管理器
	ConnMgr ziface.IConnManager
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandle.AddRouter(msgId, router)
	fmt.Println("Add router Success!")
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name is %s, listener at IP: %s, Port: %d is staring\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version is %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
		// 开启消息队列及worker工作池
		s.MsgHandle.StartWorkerPool()

		// 获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		// 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen:", s.IPVersion, ",err:", err)
			return
		}
		fmt.Println("start Zinx server success, ", s.Name, "success Listening")
		var cid uint32
		cid = 0

		// 阻塞地等待客户端连接，处理客户端链接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接的数量，则关闭此新连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too many connections! MaxConn = ", utils.GlobalObject.MaxConn)
				_ = conn.Close()
				continue		// 进行下一次的accept，直到有其它的连接业务处理完了，个数降至最大连接数量之后，这个连接再启动
			}

			// 将处理新连接的业务方法，和conn进行绑定，得到 我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandle)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
	fmt.Println("[STOP] Zinx server name is ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
	初始化Server模块的方法
 */
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		MsgHandle: NewMsgHandle(),
		ConnMgr: NewConnManager(),
	}
	return s
}