package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type Response struct {
}

func (r Response) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	response := model.Response{}
	if err := g.DecodeData(packet.Data, &response); err != nil {
		log.Println(err)
	}
	log.Println("response id", response.Id)
}
