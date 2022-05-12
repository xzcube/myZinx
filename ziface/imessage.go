package ziface

/*
	将请求的消息封装到一个Message中，定义抽象的接口
 */
type IMessage interface {
	// 获取消息的ID
	GetMsgId() uint32

	// 获取消息的长度
	GetMstLen() uint32

	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetMsgId(uint32)

	// 设置消息内容
	SetData([]byte)

	// 设置消息长度
	SetDataLen(uint32)
}