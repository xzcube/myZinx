package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myZinx/utils"
	"myZinx/ziface"
)

// DataLen (uint32 4字节) + ID(uint32 4字节)
const HEADLEN uint32 = 8

// 封包，拆包的具体模块
type DataPack struct {}

// 拆包、封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	return HEADLEN
}

// 消息封装的方法：datalen/msgId/data
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放byte字节流的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		return nil, err
	}

	// 将MsgId写入dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		return nil, err
	}

	// 将data数据，写入dataBuff中
	err = binary.Write(dataBuff, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法  将包的head信息读出来，再根据head信息里面的data的长度再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据中读取的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到dataLen和msgID
	msg := &Message{}

	// 读dataLen
	err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen)
	if err != nil {
		return nil, err
	}

	// 读MsgId
	err = binary.Read(dataBuff, binary.LittleEndian,&msg.Id)
	if err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
