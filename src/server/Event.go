package server

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/Db"
	"agent/src/server/handler"
	serverFace "agent/src/server/iface"
	serverModel "agent/src/server/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
	"sync"
	"time"
)

func NewEvent() *Event {
	e := &Event{
		methods: make(map[int32]serverFace.HandlerMethod),
	}
	e.AddHandlerMethod(0, &handler.DefaultMethod{})
	return e
}

type Event struct {
	lock    sync.RWMutex
	methods map[int32]serverFace.HandlerMethod
}

func (e *Event) AddHandlerMethod(id int32, fn serverFace.HandlerMethod) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.methods[id] = fn
}
func (e *Event) GetMethod(id int32) serverFace.HandlerMethod {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fn, ok := e.methods[id]
	if ok {
		return fn
	}
	return e.methods[0]
}
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
		db, _ := Db.DbConnections.Get("ep")
		ccServer := serverModel.Server{}
		if err := db.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
			return
		}
		auth.Id = ccServer.Id
		ccService := []*serverModel.Service{}
		if err := db.Select(&ccService, "select service_template_id as template_id,status from cc_server_service where server_id=?", ccServer.Id); err != nil {
			log.Println(err)
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
	case g.ServiceResponse:
		service := &model.ServiceResponse{}
		if err := g.DecodeData(pkt.Data, service); err != nil {
			log.Println(err)
		}
		switch service.Service {
		case g.BaseServerInfo:
			var cpu model.Cpu
			var mem model.Memory
			loadAvgs := make([]*model.LoadAvg, 0)
			var cpuNum int
			var cpuMhz string
			if err := g.DecodeData(service.Info, &cpu, &mem, &loadAvgs, &cpuNum, &cpuMhz); err != nil {
				log.Println("读取信息失败")
				break
			}
			tmp, _ := connection.GetExtraData("auth")
			auth := tmp.(*model.Auth)
			db, _ := Db.DbConnections.Get("ep")
			ram_usage_ratio := g.Round(float64(mem.Used)/float64(mem.Total), 2)
			if _, err := db.Exec("insert cc_server_log (server_id,ram,cpu_usage_ratio,ram_usage_ratio,created_at) values (?,?,?,?,?)", auth.Id, float64(mem.Total)/(1024*1024), g.Round(cpu.Busy/100, 2), ram_usage_ratio, g.CSTTime()); err != nil {
				log.Println(err.Error())
			}
			if _, err := db.Query("update cc_server set cpu_usage_ratio=?,ram_usage_ratio=?,cpu_num=?,cpu_mhz=? where id=?", g.Round(cpu.Busy/100, 2), ram_usage_ratio, cpuNum, cpuMhz, auth.Id); err != nil {
				log.Println(err.Error())
			}
		case g.HHD:
			disks := make([]*model.Disk, 0)
			if err := g.DecodeData(service.Info, &disks); err != nil {
				log.Println("读取硬盘信息失败")
				break
			}
			tmp, _ := connection.GetExtraData("auth")
			auth := tmp.(*model.Auth)
			db, _ := Db.DbConnections.Get("ep")
			for _, disk := range disks {
				var serverDisk serverModel.ServerDisk
				if err := db.Get(&serverDisk, "select id,name,gb,server_id from cc_server_disk where server_id=? and name=?", auth.Id, disk.FsFile); err != nil {
				}
				total := float64(disk.Total) / (1024 * 1024)
				if serverDisk.Id == 0 {
					if re, err := db.Exec("insert cc_server_disk (`name`,`gb`,`server_id`) values(?,?,?)", disk.FsFile, g.Round(total/1024, 2), auth.Id); err != nil {
						continue
					} else {
						serverDisk.Id, _ = re.LastInsertId()
					}
				} else {
					if _, err := db.Exec("update cc_server_disk  set gb=? where server_id=? and `name`=?", g.Round(total/1024, 2), auth.Id, disk.FsFile); err != nil {
						continue
					}
				}

				ratio := g.Round(float64(disk.Used)/float64(disk.Total), 2)
				created_at := g.CSTTime()
				if _, err := db.Exec("insert cc_server_disk_log (`disk_id`,`usage_ratio`,`created_at`) values (?,?,?)", serverDisk.Id, ratio, created_at); err != nil {
					continue
				}
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
			tmp, _ := connection.GetExtraData("auth")
			auth := tmp.(*model.Auth)
			db, _ := Db.DbConnections.Get("ep")
			created_at := g.CSTTime()
			if _, err := db.Exec("update cc_server_service set status=? where server_id=? and service_template_id=?", service.Status, auth.Id, service.Service); err != nil {
				log.Println(err)
			}
			if _, err := db.Exec("insert cc_service_log (server_service_id,status,created_at) values (?,?,?)", auth.Id, service.Status, created_at); err != nil {
				log.Println(err)
			}
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
