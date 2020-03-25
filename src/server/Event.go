package server

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	serverModel "agent/src/server/model"
	"context"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

type Event struct{}

func (e *Event) SetTimeout(connection iface.IConnection) {
	timeOut := time.Duration(utils.GlobalConfig.GetInt("heartTimeOut"))
	//每次心跳为connect设置新的过期时间,如果写入或者读取超过就会触发timeout的错误
	_ = connection.GetRawCon().SetDeadline(time.Now().Add(time.Second * timeOut))
}

func (e *Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	e.SetTimeout(connection)
}

func (e *Event) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	pkt := packet.(*src.Packet)
	switch pkt.Id {
	case g.Auth:
		var auth model.Auth
		if err := g.DecodeData(pkt.Data, &auth); err != nil {
			log.Println("读取登录信息失败,关闭连接")
			connection.Close()
			return
		}
		connection.SetExtraData("auth", &auth)
		log.Printf("agent登录,登录用户:%s\n", auth.Username)

		//用户登录成功
		pkt.Id = g.ServicesList
		db, ok := DbConnections.Get("ep")
		if !ok {
			fmt.Println("db false")
			return
		}
		fmt.Println("db ok")
		ccServer := serverModel.Server{}
		if err := db.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
			return
		}
		ccService := []*serverModel.Service{}
		if err := db.Select(&ccService, "select service_template_id as template_id,status from cc_server_service where server_id=?", ccServer.Id); err != nil {
			fmt.Println(err)
			return
		}

		service := make(map[int]int)
		for _, s := range ccService {
			service[s.TemplateId] = s.Status
		}
		fmt.Println(service)
		pkt.Data, _ = g.EncodeData(service)

		connection.Write(pkt)

	case g.PING:
		e.SetTimeout(connection)
		log.Println("心跳")
	case g.CPU:
		var cpu model.Cpu
		if err := g.DecodeData(pkt.Data, &cpu); err != nil {
			log.Println("读取cpu信息失败")
			break
		}
		log.Printf("cpu目前负载%.2f,闲置%.2f\n", cpu.Busy, cpu.Idle)
	case g.HHD:
		disks := make([]*model.Disk, 0)
		if err := g.DecodeData(pkt.Data, &disks); err != nil {
			log.Println("读取硬盘信息失败")
			break
		}
		for _, disk := range disks {
			total := float64(disk.Total) / (1024 * 1024)
			used := float64(disk.Used) / (1024 * 1024)
			free := float64(disk.Free) / (1024 * 1024)
			log.Printf("硬盘名称%s,总大小%.2fMB,已经使用%.2fMB,剩余%.2fMB\n", disk.FsFile, total, used, free)
		}
	case g.MEM:
		var mem model.Memory
		if err := g.DecodeData(pkt.Data, &mem); err != nil {
			log.Println("读取内存信息失败")
			break
		}
		log.Printf("内存大小%.2fMB,已经使用%.2fMB\n", float64(mem.Total)/(1024*1024), float64(mem.Used)/(1024*1024))
	case g.LoadAvg:
		loadAvgs := make([]*model.LoadAvg, 0)
		if err := g.DecodeData(pkt.Data, &loadAvgs); err != nil {
			log.Println("读取负载信息失败")
			break
		}
		for _, loadAvg := range loadAvgs {
			log.Printf("负载情况为%s=>%.2f\n", loadAvg.Name, loadAvg.Load)
		}
	case g.PortListen:
		ports := make([]*model.Port, 0)
		if err := g.DecodeData(pkt.Data, &ports); err != nil {
			log.Println("读取监听端口信息失败")
			break
		}
		var listenStatus string
		for _, port := range ports {
			listenStatus = "下线"
			if port.Listen {
				listenStatus = "上线"
			}
			log.Printf("监听端口协议为%s,端口号%d,监控情况%s\n", port.Type, port.Port, listenStatus)
		}
	case g.ServiceResponse:
		fmt.Println(string(pkt.Data))
	}
	packet = src.ComResponse()
	connection.Write(packet)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
	if v, ok := connection.GetExtraData("auth"); ok {
		auth := v.(*model.Auth)
		log.Printf("用户%s断开连接", auth.Username)

	}
}
