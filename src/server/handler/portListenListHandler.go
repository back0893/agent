package handler

import (
	"agent/src/g"
	"context"
	"time"

	"github.com/back0893/goTcp/iface"
)

func NewPortListenHandler() *PortListenHandler {
	return &PortListenHandler{}
}

type PortListenHandler struct {
}

func (p PortListenHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	//从mysql或者其它读取需要监听的tcp端口
	ports := []int64{80, 22, 21}
	pkt := g.NewPkt()
	pkt.Id = g.PortListen
	pkt.Data, _ = g.EncodeData(ports)
	connection.AsyncWrite(pkt, 5*time.Second)
}
