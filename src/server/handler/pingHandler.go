package handler

import (
	"agent/src"
	"context"
	"github.com/back0893/goTcp/iface"
)

type Ping struct{}

func (Ping) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
}
