package agent

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
)

type Event struct{}

func (a Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	//这个时候发送身份识别
	pkt := src.NewPkt()
	pkt.Id = g.Auth
	authModel := model.Auth{
		Username: utils.GlobalConfig.GetString("username"),
		Password: utils.GlobalConfig.GetString("password"),
	}
	pkt.Data, _ = g.EncodeData(authModel)
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
	log.Println("接连成功时")
}

func (a Event) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	pkt := packet.(*src.Packet)
	log.Println(pkt.Id)
	switch pkt.Id {

	case g.Service:
		//服务状态,服务被推动到消息队列中
		service := &model.Service{}
		if err := g.DecodeData(pkt.Data, service); err == nil {
			agent := ctx.Value(g.AGENT).(*Agent)
			agent.taskQueue.Push(service)
		} else {
			//todo 发送的消息不合规
			fmt.Println("发送的消息不合规")
		}

	case g.STOP:
		//停止
		pkt := src.NewPkt()
		pkt.Id = g.Response
		connection.Write(pkt)

		agent := ctx.Value(g.AGENT).(*Agent)
		agent.Stop()

	case g.UPDATE:
		//更新
		agent := ctx.Value(g.AGENT).(*Agent)
		info := &model.UpdateInfo{}
		_ = g.DecodeData(pkt.Data, info)
		log.Println("update", info)
		go func(agent *Agent) {
			update := NewUpdate(utils.GlobalConfig.GetString("filename"))
			if err := update.Do(info); err == nil {
				log.Println("update ok")
				agent.Stop()
			} else {
				log.Println(err)
			}
		}(agent)
		pkt := src.NewPkt()
		pkt.Id = g.Response
		connection.Write(pkt)
	case g.AuthSuccess:
		//认真成功,主动请求,启动的services
		pkt := src.NewPkt()
		pkt.Id = g.Services
		connection.Write(pkt)

	case g.AuthFail:
		//认真失败...
	case g.ServicesList:
		sl := NewServicesList()
		//如果中控服务器传递的值错误,那么就默认使用本地已经存在的services
		if err := sl.WakeUp(); err != nil {
			//读取本地保存的失败
		}
		//从中控中心获得需要启动的service
		sl.Sync(pkt.Data)

		agent := ctx.Value(g.AGENT).(*Agent)
		for _, service := range sl.GetServices() {
			agent.taskQueue.Push(service)
		}
		//保存
		sl.Sleep()
	default:
		log.Println("接受的回应id=>", pkt.Id)
	}
}

func (a Event) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("接连关闭")
}
