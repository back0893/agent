package src

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/back0893/goTcp/iface"
	"log"
	"time"
)

type Event struct{}

func (Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	connection.SetExtraData("heart", time.Now().Unix())
	log.Println("连接成功")
}

func (Event) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	pkt := packet.(*Packet)
	switch pkt.Id {
	case PING:
		connection.SetExtraData("heart", time.Now().Unix())
		log.Println("心跳")
	case CPU:
		var cpuUsage int32
		r := bytes.NewReader(pkt.Data)
		if err := binary.Read(r, binary.BigEndian, &cpuUsage); err != nil {
			log.Println("cpu使用读取失败")
		} else {
			log.Printf("cpu个数%d\n", cpuUsage)
		}
	case HHD:
		diskname := bytes.Split(pkt.Data, []byte("\n"))
		for _, name := range diskname {
			log.Printf("硬盘名称%s\n", string(name))
		}
	case MEM:
		var memTotal int64
		var memUsed int64
		r := bytes.NewReader(pkt.Data)
		if err := binary.Read(r, binary.BigEndian, &memTotal); err != nil {
			log.Println("读取失败")
			return
		}
		if err := binary.Read(r, binary.BigEndian, &memUsed); err != nil {
			log.Println("读取失败")
			return
		}
		log.Printf("内存大小%fGb,已经使用%fGb", float64(memTotal)/(1024*1024), float64(memUsed)/(1024*1024))
	}

	packet = ComResponse()
	connection.Write(packet)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("断开连接")
}
