package serviceHandler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/back0893/goTcp/iface"
)

type PortService struct {
}

func NewPortService() *PortService {
	return &PortService{}
}
func (p PortService) Handler(ctx context.Context, service *model.Service, connection iface.IConnection) error {
	ports := make([]*model.Port, 0)
	if err := g.DecodeData(service.Info, &ports); err != nil {
		log.Println("读取监听端口信息失败")
		return err
	}
	ep, ok := db.DbConnections.Get("ep")
	if !ok {
		return errors.New("db连接失败")
	}
	tmp, ok := connection.GetExtraData("auth")
	if !ok {
		return errors.New("获得用户失败")
	}
	auth := tmp.(*model.Auth)
	var listenStatus string
	for _, port := range ports {
		listenStatus = "下线"
		if port.Listen {
			listenStatus = "上线"
		}
		content := fmt.Sprintf("监听端口协议为%s,端口号%d,监控情况%s", port.Type, port.Port, listenStatus)
		if _, err := ep.Exec("insert into  cc_server_log set server_id=?,created_at=?,tag=?,content=?", auth.Id, g.CSTTime(), "server.port", content); err != nil {
			log.Println(err.Error())
		}

	}
	return nil
}
