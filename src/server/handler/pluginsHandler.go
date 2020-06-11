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
		Uri: []string{"https://github.com/m-zajac/json2go.git"},
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

type PluginsResponseHandler struct{}

func (PluginsResponseHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	metrics := make([]*model.MetricValue, 0)
	if err := g.DecodeData(packet.Data, &metrics); err != nil {
		log.Println("解析 plugins fail,response data:", string(packet.Data))
	}
	for _, metric := range metrics {
		log.Println(metric.Value)
	}
}
