package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

/**
执行命令后的返回值..
一般是处理,记录进入数据库
*/

func NewActionNotice() *ActionNotice {
	return &ActionNotice{}
}

type ActionNotice struct {
}

func (a ActionNotice) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	metrics := make([]*model.MetricValue, 0)
	if err := g.DecodeData(packet.Data, &metrics); err != nil {
		log.Println(err)
		return
	}
	for _, metric := range metrics {
		log.Println(metric.Metric, metric.Value)
	}

}
