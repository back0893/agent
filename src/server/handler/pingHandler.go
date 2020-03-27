package handler

import (
	"agent/src"
	"agent/src/server/net"
	"context"
)

type Ping struct{}

func (Ping) Handler(ctx context.Context, packet *src.Packet, connection *net.Connection) {
	connection.UpdateTimeOut()
}
