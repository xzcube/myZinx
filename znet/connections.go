package znet

import (
	"fmt"
	"myZinx/utils"
	"myZinx/ziface"
	"net"
)

/*
	连接模块
 */
type Connection struct {
	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	IsClosed bool

	// 告知当前连接已经退出的channel
	ExitChan chan bool

	// 该链接处理的方法Router
	Router ziface.IRouter
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID is", c.ConnID)
	// 启动从当前连接的读数据业务
	go c.StartReader()

	// TODO 启动从当前链接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop.. ConnID =", c.ConnID)

	// 如果当前链接已经关闭
	if c.IsClosed {
		return
	}
	c.IsClosed = false

	// 关闭socket链接
	c.Conn.Close()
	close(c.ExitChan)
}

func (c *Connection) GetConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	panic("implement me")
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		Router: router,
		IsClosed: false,  // 表示开启状态
		ExitChan: make(chan bool, 1),
	}

	return c
}

/*
	依次从用户端读取数据，读取数据之后把数据变成一个request，分别去执行用户已经重写的3个router里面的方法
 */
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")

	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buff中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buff err:", err)
			continue
		}

		// 得到当前Conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 从路由中找到注册绑定的Conn对应的router调用
		go func(request ziface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}