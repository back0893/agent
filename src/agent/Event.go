package agent

import (
	"agent/src"
	"agent/src/agent/model"
	"agent/src/g"
	"bytes"
	"context"
	"encoding/gob"
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
	if pkt.Id == g.Service {
		service := &model.Service{}
		decoder := gob.NewDecoder(bytes.NewReader(pkt.Data))
		if err := decoder.Decode(service); err == nil {
			agent := ctx.Value(g.AGENT).(*Agent)
			agent.taskQueue.Push(service)
		} else {
			//todo 发送的消息不合规
			fmt.Println("发送的消息不合规")
		}

	} else {
		log.Println("接受的回应id=>", pkt.Id)
	}
}

func (a Event) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("接连关闭")
}
