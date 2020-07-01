package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"log"

	"github.com/back0893/goTcp/iface"
)

type Response struct {
}

func (r Response) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	response := model.Common{}
	if err := g.DecodeData(packet.Data, &response); err != nil {
		log.Println(err)
	}
	log.Println("response id", response.ID)
}
