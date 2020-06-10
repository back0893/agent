package handler

import (
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type AuthFail struct {
}

func (a AuthFail) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	log.Println("认真失败")
	connection.Close()
}
