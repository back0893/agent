package g

import (
	"agent/src/g/model"
	"bytes"
	"encoding/binary"
	"github.com/back0893/goTcp/iface"
	"time"
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

/**
常用回应
*/
func ComResponse(id int32) iface.IPacket {
	pkt := NewPkt()
	pkt.Id = Response
	pkt.Data, _ = EncodeData(model.Response{Id: id})
	return pkt
}

/**
服务动作
*/
func ServicePkt(data interface{}) iface.IPacket {
	pkt := NewPkt()
	pkt.Id = 0
	pkt.Data, _ = EncodeData(data)
	return pkt
}

/**
服务的信息上报.
*/
func ServiceResponsePkt(data interface{}) iface.IPacket {
	pkt := NewPkt()
	pkt.Id = ServiceResponse
	pkt.Data, _ = EncodeData(data)
	return pkt
}
