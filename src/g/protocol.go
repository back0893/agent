package g

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/back0893/goTcp/iface"
)

var (
	PackError   = errors.New("反序列化失败")
	UnPackError = errors.New("序列化失败")
)

type Protocol struct{}

func (p Protocol) Decode(raw []byte) (iface.IPacket, error) {
	pkt := NewPkt()
	buffer := bytes.NewReader(raw)
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
	pkt.Data = make([]byte, length)
	if err := binary.Read(buffer, binary.BigEndian, pkt.Data); err != nil {
		return nil, UnPackError
	}
	return pkt, nil
}

func (Protocol) Pack(pack iface.IPacket) ([]byte, error) {
	raw, err := pack.Serialize()
	if err != nil {
		return nil, PackError
	}
	return raw, nil
}

func (Protocol) UnPack(data []byte, atEOF bool) (advance int, token []byte, err error) {
	dataLength := len(data)
	//如果长度超过长度,应该是解析错误.
	if dataLength > 65535 {
		return 0, nil, errors.New("长度过长")
	}
	if dataLength < 8 {
		return 0, nil, nil
	}

	buffer := bytes.NewBuffer(data[:8])
	var length int64
	if err := binary.Read(buffer, binary.BigEndian, &length); err != nil {
		return 0, nil, UnPackError
	}
	if dataLength < int(length) {
		if atEOF {
			return 0, nil, UnPackError
		}
		return 0, nil, nil
	}

	return int(length), data[:length], nil
}
