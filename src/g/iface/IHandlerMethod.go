package iface

import (
	"agent/src"
	"context"
	"github.com/back0893/goTcp/iface"
)

type IHandlerMethod interface {
	Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection)
}
