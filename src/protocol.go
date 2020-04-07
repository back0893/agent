package src

import (
	"agent/src/g"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"log"
	"time"
)

var (
	PackError   = errors.New("反序列化失败")
	UnPackError = errors.New("序列化失败")
)

/**
出现了.因为错误导致的读取包出现错误的问题.
还是新增一个开始和结束分割符号
学习808的协议使用7e 作为开始和结束的标示
负载中的 7e =>7d02
负载中的 7d =>7d01
*/
type Protocol struct{}

func (Protocol) Pack(pack iface.IPacket) ([]byte, error) {
	raw, err := pack.Serialize()
	raw = bytes.ReplaceAll(raw, []byte{0x7e}, []byte{0x7d, 0x02})
	raw = bytes.ReplaceAll(raw, []byte{0x7d}, []byte{0x7d, 0x01})
	raw = append([]byte{0x7e}, raw...)
	raw = append(raw, 0x7e)
	if err != nil {
		return nil, PackError
	}
	return raw, nil
}

func (Protocol) UnPack(conn iface.IConnection) (iface.IPacket, error) {
	pkt := NewPkt()
	buffer := conn.GetBuffer()
	_, err := buffer.ReadBytes(0x7e)
	if err != nil {
		return nil, err
	}
	read := make(chan struct{})
	var data []byte
	go func() {
		data, err = buffer.ReadBytes(0x7e)
		if err != nil {
			return
		}
		read <- struct{}{}
	}()
	//一个简单的读取定时器
	tick := time.NewTicker(20 * time.Second)
	select {
	case <-tick.C:
		fmt.Println("超时")
		return nil, UnPackError
	case <-read:
		tick.Stop()
	}
	if len(data) > 1048576 {
		return nil, UnPackError
	}
	data = bytes.ReplaceAll(data, []byte{0x7d, 0x02}, []byte{0x7e})
	data = bytes.ReplaceAll(data, []byte{0x7d, 0x01}, []byte{0x7d})
	if len(data) < g.HeaderLength {
		return nil, UnPackError
	}
	b := bytes.NewReader(data)
	if err := binary.Read(b, binary.BigEndian, &pkt.Length); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(b, binary.BigEndian, &pkt.Version); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(b, binary.BigEndian, &pkt.Timestamp); err != nil {
		return nil, UnPackError
	}
	if err := binary.Read(b, binary.BigEndian, &pkt.Id); err != nil {
		return nil, UnPackError
	}

	//固定负载数据
	length := pkt.Length - g.HeaderLength
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
