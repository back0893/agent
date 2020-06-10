package g

import (
	"encoding/binary"
	"errors"
	"github.com/back0893/goTcp/iface"
	"log"
)

var (
	PackError   = errors.New("反序列化失败")
	UnPackError = errors.New("序列化失败")
)

type Protocol struct{}

func (Protocol) Pack(pack iface.IPacket) ([]byte, error) {
	raw, err := pack.Serialize()
	if err != nil {
		return nil, PackError
	}
	return raw, nil
}

func (Protocol) UnPack(conn iface.IConnection) (iface.IPacket, error) {
	pkt := NewPkt()
	buffer := conn.GetBuffer()
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Length); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Version); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Timestamp); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Id); err != nil {
		return nil, UnPackError
	}

	//固定负载数据
	length := pkt.Length - HeaderLength
	if length < 0 {
		log.Println("长度不足", pkt.Id, pkt.Length, pkt.Timestamp)
		return nil, errors.New("长度不足")
	}
	//如果长度超过长度,应该是解析错误.
	if length > 1024*100 {
		log.Println("长度过长", pkt.Id, pkt.Length, pkt.Timestamp)
		return nil, errors.New("长度过长")
	}
	pkt.Data = make([]byte, length)
	if err := binary.Read(buffer, binary.BigEndian, pkt.Data); err != nil {
		return nil, UnPackError
	}
	return pkt, nil
}
