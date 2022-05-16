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
	// 当前的Server添加又给router，server注册链接对应的业务
	Router ziface.IRouter
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router Success!")
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name is %s, listener at IP: %s, Port: %d is staring\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx] Version is %s, MaxConn: %d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	go func() {
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

			// 将处理新连接的业务方法，和conn进行绑定，得到 我们的链接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的连接信息 进行停止或者回收
}

func (s *Server) Server() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞状态
	select {

	}
}

/*
	初始化Server模块的方法
 */
func NewServer() ziface.IServer {
	s := &Server{
		Name: utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP: utils.GlobalObject.Host,
		Port: utils.GlobalObject.TcpPort,
		Router: nil,
	}
	return s
}