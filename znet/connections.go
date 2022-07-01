package znet

import (
	"errors"
	"fmt"
	"io"
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

	// 消息的管理msgId 和对应的处理业务api
	MsgHandle ziface.IMsgHandle
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
	_ = c.Conn.Close()
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

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed {
		return errors.New("Connection closed ")
	}

	// 将data进行封包 MsgDataLen + MsgID + MsgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg error, id is:", msgId)
		return errors.New("Pack error msg ")
	}

	// 将数据发送给客户端
	_, err = c.Conn.Write(binaryMsg)
	if err != nil {
		fmt.Println("Write msg err, msgId is:", msgId)
		return errors.New("conn Write error")
	}

	return nil
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, handle ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn: conn,
		ConnID: connID,
		MsgHandle: handle,
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
		/*buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buff err:", err)
			continue
		}*/

		// 创建一个拆包解包的对象
		dp := NewDataPack()

		// 读取客户端的MsgHead的二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetConnection(), headData); err != nil {
			fmt.Println("read msgHead err:", err)
		}

		// 拆包，得到msgID和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
		}

		// 根据dataLen 再次读取data，放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetConnection(), data)
			if err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
		}

		msg.SetData(data)

		// 得到当前Conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 从路由中找到注册绑定的conn对应的router调用
		// 根据绑定好的MsgID,找到对应处理api的业务并执行
		go c.MsgHandle.DoMsgHandle(&req)
	}
}