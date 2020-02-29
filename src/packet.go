package src

import (
	"bytes"
	"encoding/binary"
	"github.com/back0893/goTcp/iface"
	"time"
)

const (
	HeaderLength = 21  //包的固定长度
	PING         = 1   //心跳
	Auth         = 100 //身份识别
	CPU          = 101 //cpu使用
	HHD          = 102 //硬盘使用
	MEM          = 103 //内存使用
	LoadAvg      = 104 //负载
	PortListen   = 105 //端口监听情况
	Response     = 100 //通用回应
)

type Packet struct {
	Length    int64  //数据长度
	Version   int8   //版本号
	Timestamp int64  //发送时间
	Id        int32  //命令id
	Data      []byte //负载数据

}

func NewPkt() *Packet {
	p := Packet{
		Timestamp: time.Now().Unix(),
		Data:      []byte{},
		Version:   1,
	}
	return &p
}

func (pkt *Packet) Len() int {
	var length int = HeaderLength
	length += len(pkt.Data)
	return length
}
func (pkt *Packet) Serialize() ([]byte, error) {
	var buffer bytes.Buffer

	pkt.Length = int64(pkt.Len())

	if err := binary.Write(&buffer, binary.BigEndian, pkt.Length); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, pkt.Version); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, pkt.Timestamp); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, pkt.Id); err != nil {
		return nil, err
	}
	if err := binary.Write(&buffer, binary.BigEndian, pkt.Data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func ComResponse() iface.IPacket {
	pkt := NewPkt()
	pkt.Id = Response
	return pkt
}
