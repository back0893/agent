package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"context"
	"log"

	"github.com/back0893/goTcp/iface"
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
	tmp, ok := connection.GetExtraData("auth")
	if !ok {
		return
	}
	auth := tmp.(*model.Auth)
	ep, ok := db.DbConnections.Get("ep")
	if !ok {
		return
	}
	for _, metric := range metrics {
		log.Println(metric.Metric, metric.Value)
		if _, err := ep.Exec("insert into cc_server_log set server_id=?,created_at=?,tag=?,content=?", auth.Id, g.CSTTime(), metric.Metric, metric.Value); err != nil {
			log.Println(err.Error())
		}
	}
}
