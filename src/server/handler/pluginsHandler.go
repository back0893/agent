package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"log"
)

func NewPluginsHandler() *PluginsHandler {
	return &PluginsHandler{}
}

type PluginsHandler struct{}

func (PluginsHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {

	pkt := g.NewPkt()
	plugin := model.Plugins{
		Uri: []string{"git@gitee.com:liuzy1988/SealKit.git"},
	}
	var err error
	if pkt.Data, err = g.EncodeData(plugin); err != nil {
		log.Println(err)
		return
	}
	fmt.Println(len(pkt.Data))
	pkt.Id = g.MinePluginsResponse
	connection.Write(pkt)
}
