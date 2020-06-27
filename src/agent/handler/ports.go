package handler

import (
	g2 "agent/src/agent/g"
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type Ports struct {
}

func (p Ports) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	ports := make([]int64, 0)
	if err := g.DecodeData(packet.Data, &ports); err != nil {
		log.Println(err)
		return
	}
	g2.SetPortListen(ports)
}
