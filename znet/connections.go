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
	panic("implement me")
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop.. ConnID =", c.ConnID)
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

