package handler

import (
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type BackDoorHandler struct {
}

func NewBackDoorHandler() *BackDoorHandler {
	return &BackDoorHandler{}
}

func (b BackDoorHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	response := ""
	if err := g.DecodeData(packet.Data, &response); err != nil {
		log.Println(err)
	}
	log.Println(response)
}
