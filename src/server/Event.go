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

		log.Printf("agent登录,登录用户:%s\n", auth.Username)
		db, _ := DbConnections.Get("ep")
		ccServer := serverModel.Server{}
		if err := db.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
			return
		}
		auth.Id = ccServer.Id
		ccService := []*serverModel.Service{}
		if err := db.Select(&ccService, "select service_template_id as template_id,status from cc_server_service where server_id=?", ccServer.Id); err != nil {
			fmt.Println(err)
			return
		}
		connection.SetExtraData("auth", &auth)

		service := make(map[int]int)
		for _, s := range ccService {
			service[s.TemplateId] = s.Status
		}
		//用户登录成功
		pkt.Id = g.ServicesList

		pkt.Data, _ = g.EncodeData(service)

		connection.Write(pkt)

	case g.PING:
		e.SetTimeout(connection)
		log.Println("心跳")
	case g.ServiceResponse:
		service := &model.ServiceResponse{}
		if err := g.DecodeData(pkt.Data, service); err != nil {
			fmt.Println(err)
		}
		switch service.Service {
		case g.BaseServerInfo:
			var cpu model.Cpu
			var mem model.Memory
			loadAvgs := make([]*model.LoadAvg, 0)
			if err := g.DecodeData(service.Info, &cpu, &mem, &loadAvgs); err != nil {
				log.Println("读取信息失败")
				break
			}
			tmp, _ := connection.GetExtraData("auth")
			auth := tmp.(*model.Auth)
			db, _ := DbConnections.Get("ep")
			memBusy := (mem.Used * 10000 / mem.Total) / 100
			if _, err := db.Exec("insert cc_server_log (server_id,ram,cpu_usage_ratio,ram_usage_ratio) values (?,?,?,?)", auth.Id, float64(mem.Total)/(1024*1024), cpu.Busy, memBusy); err != nil {
				fmt.Println(err.Error())
			}
			if _, err := db.Query("update cc_server set cpu_usage_ratio=?,ram_usage_ratio=? where id=?", cpu.Busy, memBusy, auth.Id); err != nil {
				fmt.Println(err.Error())
			}
			log.Printf("cpu目前负载%.2f,闲置%.2f\n", cpu.Busy, cpu.Idle)
			log.Printf("内存大小%.2fMB,已经使用%.2fMB\n", float64(mem.Total)/(1024*1024), float64(mem.Used)/(1024*1024))
		case g.HHD:
			disks := make([]*model.Disk, 0)
			if err := g.DecodeData(service.Info, &disks); err != nil {
				log.Println("读取硬盘信息失败")
				break
			}
			for _, disk := range disks {
				total := float64(disk.Total) / (1024 * 1024)
				used := float64(disk.Used) / (1024 * 1024)
				free := float64(disk.Free) / (1024 * 1024)
				log.Printf("硬盘名称%s,总大小%.2fMB,已经使用%.2fMB,剩余%.2fMB\n", disk.FsFile, total, used, free)
			}
		case g.PortListen:
			ports := make([]*model.Port, 0)
			if err := g.DecodeData(service.Info, &ports); err != nil {
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
		case g.REDISSERVICE:
			var info string
			if err := g.DecodeData(service.Info, &info); err != nil {
				log.Println("读取redis失败")
				break
			}
			log.Println("redis====>", info)
		}
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
