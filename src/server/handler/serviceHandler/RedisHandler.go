package serviceHandler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/Db"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type RedisService struct {
}

func (r RedisService) Handler(ctx context.Context, service *model.ServiceResponse, connection iface.IConnection) error {
	var info string
	if err := g.DecodeData(service.Info, &info); err != nil {
		log.Println("读取redis失败")
		return err
	}
	tmp, _ := connection.GetExtraData("auth")
	auth := tmp.(*model.Auth)
	db, _ := Db.DbConnections.Get("ep")
	created_at := g.CSTTime()
	if _, err := db.Exec("update cc_server_service set status=? where server_id=? and service_template_id=?", service.Status, auth.Id, service.Service); err != nil {
		log.Println(err)
		return err
	}
	if _, err := db.Exec("insert cc_service_log (server_service_id,status,created_at) values (?,?,?)", auth.Id, service.Status, created_at); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
