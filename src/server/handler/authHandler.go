package handler

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/Db"
	serverModel "agent/src/server/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

type AuthHandler struct{}

func (AuthHandler) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	var auth model.Auth
	if err := g.DecodeData(packet.Data, &auth); err != nil {
		log.Println("读取登录信息失败,关闭连接")
		connection.Close()
		return
	}

	log.Printf("agent登录,登录用户:%s\n", auth.Username)
	db, _ := Db.DbConnections.Get("ep")
	ccServer := serverModel.Server{}
	if err := db.Get(&ccServer, "select id,name from cc_server where name=?", auth.Username); err != nil {
		return
	}
	auth.Id = ccServer.Id
	ccService := []*serverModel.Service{}
	if err := db.Select(&ccService, "select service_template_id as template_id,status from cc_server_service where server_id=?", ccServer.Id); err != nil {
		log.Println(err)
		return
	}
	connection.SetExtraData("auth", &auth)

	service := make(map[int]int)
	for _, s := range ccService {
		service[s.TemplateId] = s.Status
	}
	//用户登录成功否下发服务
	packet.Id = g.ServicesList

	packet.Data, _ = g.EncodeData(service)

	connection.Write(packet)
}
