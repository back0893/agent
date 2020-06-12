package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"time"
)

type Update struct {
}

func (u Update) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	info := model.UpdateInfo{}

	//回应接收到更新通知
	pkt := g.ComResponse(packet.Id)
	connection.AsyncWrite(pkt, 5*time.Second)
	if err := g.DecodeData(packet.Data, &info); err == nil {
	}
	ch := ctx.Value("upgradeChan").(chan *model.UpdateInfo)
	ch <- &info
}
