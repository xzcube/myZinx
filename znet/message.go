package znet

type Message struct {
	Id      uint32 // 消息的ID
	DataLen uint32 // 消息的长度
	Data    []byte // 消息的数据
}

func (m *Message) GetMsgId() uint32 {
	panic("implement me")
}

func (m *Message) GetMstLen() uint32 {
	panic("implement me")
}

func (m *Message) GetData() []byte {
	panic("implement me")
}

func (m *Message) SetMsgId(id uint32) {
	panic("implement me")
}

func (m *Message) SetData(data []byte) {
	panic("implement me")
}

func (m *Message) SetDataLen(len uint32) {
	panic("implement me")
}
