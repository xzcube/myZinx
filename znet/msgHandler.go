package znet

import (
	"fmt"
	"myZinx/ziface"
	"strconv"
)


/*
	消息处理模块的实现
 */
type MsgHandle struct{
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32] ziface.IRouter
}

// 提供一个初始化MsgHandle的方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32] ziface.IRouter),
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

