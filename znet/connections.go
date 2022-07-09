package znet

import (
	"errors"
	"fmt"
	"io"
	"myZinx/utils"
	"myZinx/ziface"
	"net"
)

/*
	连接模块
 */
type Connection struct {
	// 当前connection是隶属于哪个Server的
	TCPServer ziface.IServer

	// 当前连接的socket TCP套接字
	Conn *net.TCPConn

	// 连接的ID
	ConnID uint32

	// 当前的连接状态
	IsClosed bool

	// 告知当前连接已经退出的channel(由Reader告知Writer退出)
	ExitChan chan bool

	// 无缓冲的管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理msgId 和对应的处理业务api
	MsgHandle ziface.IMsgHandle
}

// 初始化连接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, handle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TCPServer: server,
		Conn: conn,
		ConnID: connID,
		MsgHandle: handle,
		IsClosed: false,  // 表示开启状态
		ExitChan: make(chan bool, 1),
		msgChan: make(chan []byte),
	}

	// 将conn加入到ConnManager中
	c.TCPServer.GetConnMgr().Add(c)

	return c
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID is", c.ConnID)
	// 启动从当前连接的读数据业务
	go c.StartReader()

	// 启动从当前链接写数据的业务
	go c.StartWriter()

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

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前连接从connection中摘除掉
	c.TCPServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)

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
	c.msgChan <- binaryMsg

	return nil
}

/*
	依次从用户端读取数据，读取数据之后把数据变成一个request，分别去执行用户已经重写的3个router里面的方法
 */
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")

	defer fmt.Println("[Reader is exit],connID = ", c.ConnID, " remote addr is", c.GetRemoteAddr().String())
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
			break
		}

		// 拆包，得到msgID和msgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:", err)
			break
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

		// 做一个判断，判断是否已经开启了工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册绑定的conn对应的router调用
			// 根据绑定好的MsgID,找到对应处理api的业务并执行
			go c.MsgHandle.DoMsgHandle(&req)
		}

	}
}

/*
	写消息的goroutine，专门发送给客户端消息的模块
 */
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[conn Writer exit!] ", c.GetRemoteAddr().String())

	// 不断地阻塞等待channel消息，写给客户端
	for {
		select {
		case data := <- c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("send data error:", err)
				return
			}
		case <- c.ExitChan:
			// 代表reader已经退出，此时Writer也要退出
			return
		}
	}
}