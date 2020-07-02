package handler

import (
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
)

type Stop struct{}

//Handler 停止处理
func(s Stop)Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection){
	connection.Close()
}