package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"context"
	"log"
	"strings"

	"github.com/back0893/goTcp/iface"
)

type UpdateHandler struct{}

func (u UpdateHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	info := model.UpdateResponse{}
	if err := g.DecodeData(packet.Data, &info); err != nil {
		log.Println(err)
		return
	}
	ep, ok := db.DbConnections.Get("ep")
	if !ok {
		return
	}
	status := 1
	if info.Status {
		status = 2
	}
	if _, err := ep.Exec("update cc_service_log set status=?,log=? where id=?", status, strings.Trim(info.Message, "\n\r"), info.LogID); err != nil {
		log.Println(err.Error())
		return
	}
}
