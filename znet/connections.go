package znet

import (
	"fmt"
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

	// 当前连接所绑定的处理业务的方法API
	HandleAPI ziface.HandleFunc

	// 告知当前连接已经退出的channel
	ExitChan chan bool
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
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		HandleAPI: callbackApi,
		IsClosed: false,  // 表示开启状态
		ExitChan: make(chan bool, 1),
	}

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")

	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is", c.GetRemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buff中，目前最大为512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buff err:", err)
			continue
		}

		// 调用当前连接所绑定的handleAPI
		err = c.HandleAPI(c.Conn, buf, cnt)
		if err != nil {
			fmt.Println("ConnID:", c.ConnID, "handle is error:", err)
			break
		}
	}

}