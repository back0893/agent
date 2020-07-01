package serviceHandler

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
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
	var listenStatus string
	for _, port := range ports {
		listenStatus = "下线"
		if port.Listen {
			listenStatus = "上线"
		}
		log.Printf("监听端口协议为%s,端口号%d,监控情况%s\n", port.Type, port.Port, listenStatus)
	}
	return nil
}
