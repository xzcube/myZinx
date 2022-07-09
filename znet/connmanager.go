package znet

import (
	"errors"
	"fmt"
	"myZinx/ziface"
	"sync"
)

type ConnManger struct {
	connections map[uint32]ziface.IConnections // 管理的连接信息集合
	connLock    sync.RWMutex                   // 保护连接集合的读写锁
}

// 创建当前连接的方法
func NewConnManager() *ConnManger {
	return &ConnManger{
		connections: make(map[uint32] ziface.IConnections),
	}
}

func (connMgr *ConnManger) Add(conn ziface.IConnections) {
	// 保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock() 		// 最后要把锁打开

	// 将conn加入到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully:conn num = ", connMgr.Len())
}

func (connMgr *ConnManger) Remove(conn ziface.IConnections) {
	// 保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("remove from ConnManager successfully:conn num = ", connMgr.Len())
}

func (connMgr *ConnManger) Get(connID uint32) (ziface.IConnections, error) {
	// 保护共享资源map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		// 找到了
		return conn, nil
	}else {
		return nil, errors.New("connection not found, connId is")
	}
}

func (connMgr *ConnManger) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManger) ClearConn() {
	// 保护共享资源map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}

	fmt.Println("Clear All connections success!")
}
