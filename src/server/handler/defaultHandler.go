package handler

import (
	"agent/src"
	"agent/src/server/net"
	"context"
	"log"
)

type DefaultMethod struct {
}

func (d DefaultMethod) Handler(ctx context.Context, packet *src.Packet, connection *net.Connection) {
	log.Printf("方法还未被时实现method_id===>%d", packet.Id)
	pkt := src.ComResponse()
	connection.Write(pkt)
}
