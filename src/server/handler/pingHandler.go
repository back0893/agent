package handler

import (
	"agent/src"
	"context"
	"github.com/back0893/goTcp/iface"
)

func NewPing() *Ping {
	return &Ping{}
}

type Ping struct{}

func (Ping) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	//todo 心跳需要做的事情
	pkt := src.ComResponse()
	connection.Write(pkt)
}
