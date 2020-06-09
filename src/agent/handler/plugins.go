package handler

import (
	"agent/src"
	"context"
	"github.com/back0893/goTcp/iface"
)

type Plugins struct {
}

func (p Plugins) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {

}
