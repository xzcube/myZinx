package znet

import (
	"fmt"
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
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP is %s, Port is %d. is staring\n", s.IP, s.Port)

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

		// 阻塞地等待客户端连接，处理客户端链接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}

			// 已经与客户端建立连接，做一些业务，做一个最基本的最大512字节的回显业务
			go func() {
				for {
					buff := make([]byte, 512)
					cnt, err := conn.Read(buff)
					if err != nil {
						fmt.Println("recv buf err, ", err)
						continue
					}

					// 回显功能
					if _, err := conn.Write(buff[:cnt]); err != nil {
						// 回显失败
						fmt.Println("write back buff err:", err)
						continue
					}
				}

			}()
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
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name: name,
		IPVersion: "tcp4",
		IP: "0.0.0.0",
		Port: 8999,
	}
	return s
}