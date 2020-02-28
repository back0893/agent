package src

import (
	"agent/src/agent/model"
	"bytes"
	"context"
	"encoding/gob"
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
	r := bytes.NewReader(pkt.Data)
	decoder := gob.NewDecoder(r)
	switch pkt.Id {
	case PING:
		connection.SetExtraData("heart", time.Now().Unix())
		log.Println("心跳")
	case CPU:
		var cpu model.Cpu
		if err := decoder.Decode(&cpu); err != nil {
			log.Println("读取cpu信息失败")
			break
		}
		log.Printf("cpu目前闲置%.2f,负载%.2f\n", cpu.Busy, cpu.Idle)
	case HHD:
		disks := make([]*model.Disk, 0)
		if err := decoder.Decode(&disks); err != nil {
			log.Println("读取硬盘信息失败")
			break
		}
		for _, disk := range disks {
			total := float64(disk.Total) / (1024 * 1024)
			used := float64(disk.Used) / (1024 * 1024)
			free := float64(disk.Free) / (1024 * 1024)
			log.Printf("硬盘名称%s\n,总大小%.2fMB,已经使用%.2fMB,剩余%.2fMB\n", disk.FsFile, total, used, free)
		}
	case MEM:
		var mem model.Memory
		if err := decoder.Decode(&mem); err != nil {
			log.Println("读取内存信息失败")
			break
		}
		log.Printf("内存大小%.2fMB,已经使用%.2fMB\n", float64(mem.Total)/(1024*1024), float64(mem.Used)/(1024*1024))
	case LoadAvg:
		loadAvgs := make([]*model.LoadAvg, 0)
		if err := decoder.Decode(&loadAvgs); err != nil {
			log.Println("读取负载信息失败")
			break
		}
		for _, loadAvg := range loadAvgs {
			log.Printf("负载情况为%s=>%.2f\n", loadAvg.Name, loadAvg.Load)
		}
	case PortListen:
		ports := make([]*model.Port, 0)
		if err := decoder.Decode(&ports); err != nil {
			log.Println("读取监听端口信息失败")
			break
		}
		for _, port := range ports {
			log.Printf("监听端口协议为%s,端口号%d\n", port.Type, port.Port)
		}
	}
	packet = ComResponse()
	connection.Write(packet)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("断开连接")
}
