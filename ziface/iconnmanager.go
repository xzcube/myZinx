package ziface

/*
	连接管理模块抽象层
	创建一个连接管理模块	将连接管理模块集中到Zinx框架中
	给Zinx框架提供创建连接之后/销毁连接之前所要处理的一些业务  例如提供给用户能够注册的Hook函数
 */
type IConnManager interface {
	// 添加连接
	Add(conn IConnections)

	// 删除连接
	Remove(conn IConnections)

	// 根据connID获取连接
	Get(connID uint32) (IConnections, error)

	// 得到当前连接总数
	Len() int

	// 清除并终止所有连接
	ClearConn()
}
