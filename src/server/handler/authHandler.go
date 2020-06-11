package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	serverModel "agent/src/server/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
	"time"
)

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type AuthHandler struct{}

func (AuthHandler) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	var auth model.Auth
	if err := g.DecodeData(packet.Data, &auth); err != nil {
		log.Println("读取登录信息失败,关闭连接")
		connection.Close()
		return
	}

	log.Printf("agent登录,登录用户:%s\n", auth.Username)
	ep, _ := db.DbConnections.Get("ep")
	ccServer := serverModel.Server{}

	pkt := g.NewPkt()
	pkt.Id = g.AuthSuccess
	if err := ep.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
		pkt.Id = g.AuthFail
		connection.AsyncWrite(pkt, 5*time.Second)
		return
	}
	auth.Id = ccServer.Id
	ccService := []*serverModel.Service{}
	if err := ep.Select(&ccService, "select service_template_id as template_id,status from cc_server_service where server_id=?", ccServer.Id); err != nil {
		log.Println(err)
		pkt.Id = g.AuthFail
		connection.AsyncWrite(pkt, 5*time.Second)
		return
	}
	connection.AsyncWrite(pkt, 5*time.Second)

	connection.SetExtraData("auth", &auth)

	service := make(map[int]int)
	for _, s := range ccService {
		service[s.TemplateId] = s.Status
	}
}
