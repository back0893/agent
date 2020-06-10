package agent

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
)

type Event struct {
	*g.Event
}

func NewEvent() *Event {
	return &Event{
		g.NewEvent(),
	}
}

func (a Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	//这个时候发送身份识别
	pkt := g.NewPkt()
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

func (a Event) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("接连关闭")
}
