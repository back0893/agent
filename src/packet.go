package src

import (
	"agent/src/agent/model"
	"agent/src/g"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
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
	var length int = g.HeaderLength
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
func ComResponse() iface.IPacket {
	pkt := NewPkt()
	pkt.Id = g.Response
	return pkt
}
func ServicePkt(service *model.Service) iface.IPacket {
	pkt := NewPkt()
	pkt.Id = g.Service
	fmt.Println(service)
	buf := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buf)
	_ = encoder.Encode(service)
	pkt.Data = buf.Bytes()
	return pkt
}
