package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	serverModel "agent/src/server/model"
	"context"
	"log"
	"time"

	"github.com/back0893/goTcp/iface"
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
	pkt.Id = g.Auth
	data := model.AuthResponse{}
	if err := ep.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
		connection.AsyncWrite(pkt, 5*time.Second)
		pkt.Data, _ = g.EncodeData(data)
		return
	}
	auth.Id = ccServer.Id

	data.Status = true
	pkt.Data, _ = g.EncodeData(data)
	connection.AsyncWrite(pkt, 5*time.Second)

	connection.SetExtraData("auth", &auth)
}
