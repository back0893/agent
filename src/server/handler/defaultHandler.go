package handler

import (
	"agent/src"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type DefaultMethod struct {
}

func NewDefaultMethod() *DefaultMethod {
	return &DefaultMethod{}
}

func (d DefaultMethod) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	log.Printf("方法还未被时实现method_id===>%d", packet.Id)
	pkt := src.ComResponse()
	connection.Write(pkt)
}
