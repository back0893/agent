package iface

import (
	"agent/src"
	"agent/src/server/net"
	"context"
)

type HandlerMethod interface {
	Handler(ctx context.Context, packet *src.Packet, connection *net.Connection)
}
