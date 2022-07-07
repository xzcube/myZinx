package znet

import (
	"fmt"
	"myZinx/utils"
	"myZinx/ziface"
	"strconv"
)


/*
	消息处理模块的实现
 */
type MsgHandle struct{
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32] ziface.IRouter

	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// 提供一个初始化MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,	// 从全局配置中获取
		TaskQueue: make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	// 从request中找到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		// 没有这个方法
		fmt.Println("api msgID = ", request.GetMsgID(), "is NOT FOUND! need register")
	}

	// 存在这个方法，根据msgID直接调用即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)

}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msg绑定的api处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// 当前msgID对应的api是存在的
		panic("repeat api, msgID is" + strconv.Itoa(int(msgID)))
	}

	// 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("add api MsgID = ", msgID, "success!")
}

// 启动一个worker工作池(开启工作池的动作只能发生一次，因为一个zinx框架只能有一个工作池)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize 分别开启worker，每个worker用一个go程来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1.给当前worker对应的channel消息队列开辟管道空间 第零个worker就用第零个channel，第i个worker就用第i个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)

		// 2.启动当前的worker，阻塞等待消息从channel中传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(workId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID = ", workId, " is started..")

	// 不断地阻塞等待对应的消息队列
	for {
		select {
		case request := <-taskQueue:
			// 如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务
			mh.DoMsgHandle(request)
		}
	}
}

// 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1.将消息平均分配给不同的worker
	// 根据客户端建立的ConnID进行分配
	// 基本的平均分配的轮询法则
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(), " to workerID = ", workerID)

	// 2.将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}