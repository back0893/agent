package handler

import (
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type AuthSuccess struct {
}

func (a AuthSuccess) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	//同步监控的进程和端口
	//发送通知,通知中控服务器下发监控的端口号
	pkt := g.NewPkt()
	pkt.Id = g.PortListenList
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}

	//发送通知,通知中控服务器下发监控的进程id
	pkt = g.NewPkt()
	pkt.Id = g.ProcessNumList
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
	//同步插件
	//发送通知,通知中控服务器下发插件
	pkt = g.NewPkt()
	pkt.Id = g.MinePlugins
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
}
