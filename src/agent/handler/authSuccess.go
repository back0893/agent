package handler

import (
	"agent/src"
	"agent/src/agent/cron"
	"agent/src/agent/funcs"
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type AuthSuccess struct {
}

func (a AuthSuccess) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	//认真成功,主动请求,启动的services
	//初始化服务器的数据收集
	funcs.BuildMappers()
	//定时更新cpu的使用情况
	cron.InitDatHistory()

	//同步监控的进程和端口
	//发送通知,通知中控服务器下发监控的端口号
	pkt := src.NewPkt()
	pkt.Id = g.PortListenList
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}

	//发送通知,通知中控服务器下发监控的进程id
	pkt.Id = g.ProcessNumList
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
	//同步插件
	//发送通知,通知中控服务器下发插件
	pkt.Id = g.MinePlugins
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
}
