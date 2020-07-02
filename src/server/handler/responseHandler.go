package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"context"
	"log"

	"github.com/back0893/goTcp/iface"
)

func NewResponseHandler() *ResponseHandler {
	return &ResponseHandler{}
}

type ResponseHandler struct{}

func (p ResponseHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	info := model.Common{}
	if err := g.DecodeData(packet.Data, &info); err != nil {
		log.Println("读取失败", packet.Id)
		return
	}
	ep, ok := db.DbConnections.Get("ep")
	if !ok {
		log.Println("db连接失败")
		return
	}
	status := 1
	if info.Status == 0 {
		status = 2
	}
	if _, err := ep.Exec("update cc_service_log set status=?,log=? where id=?", status, info.Message, info.ID); err != nil {
		log.Println(err.Error())
		return
	}
	return
}
