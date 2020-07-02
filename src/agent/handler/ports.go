package handler

import (
	g2 "agent/src/agent/g"
	"agent/src/g"
	"context"
	"log"

	"github.com/back0893/goTcp/iface"
)

type Ports struct {
}

func (p Ports) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	ports := make([]int64, 0)
	var logID int32
	if err := g.DecodeData(packet.Data, &ports, &logID); err != nil {
		log.Println(err)
		return
	}
	g2.SetPortListen(ports)
	pkt := g.ComResponse(logID, 1, "接收成功")
	connection.Write(pkt)
}
