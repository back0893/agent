package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"

	"github.com/back0893/goTcp/iface"
)

type Update struct {
}

func (u Update) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	info := model.UpdateInfo{}
	var logID int32
	if err := g.DecodeData(packet.Data, &info, &logID); err == nil {
	}
	info.LogID = logID
	ch := ctx.Value("upgradeChan").(chan *model.UpdateInfo)
	//回应接收到更新通知
	ch <- &info

}
