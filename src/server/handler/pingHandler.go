package handler

import (
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"time"
)

func NewPing() *Ping {
	return &Ping{}
}

type Ping struct{}

func (Ping) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	connection.SetExtraData("last_ping", time.Now().Unix())
	//todo 心跳需要做的事情
	pkt := g.ComResponse()
	connection.Write(pkt)
}
