package handler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/handler/serviceHandler"
	serverFace "agent/src/server/iface"
	"context"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"log"
)

type ServiceResponse struct {
	methods map[int]serverFace.ServiceMethod
}

func NewServiceResponse() *ServiceResponse {
	sr := &ServiceResponse{
		methods: make(map[int]serverFace.ServiceMethod),
	}
	sr.AddHandlerMethod(g.BaseServerInfo, serviceHandler.NewBaseServerInfo())
	sr.AddHandlerMethod(g.HHD, serviceHandler.NewHHDHandler())
	sr.AddHandlerMethod(g.PortListen, serviceHandler.NewPortService())
	return sr
}
func (sr *ServiceResponse) AddHandlerMethod(id int, fn serverFace.ServiceMethod) {
	sr.methods[id] = fn
}
func (sr *ServiceResponse) GetMethod(id int) serverFace.ServiceMethod {
	fn, ok := sr.methods[id]
	if ok {
		return fn
	}
	return nil
}

func (sr *ServiceResponse) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	service := &model.ServiceResponse{}
	if err := g.DecodeData(packet.Data, service); err != nil {
		log.Println(err)
	}
	id := service.Service
	fn := sr.GetMethod(id)
	if fn == nil {
		//todo 没有不做任何处理
		return
	}
	if err := fn.Handler(ctx, service, connection); err != nil {
		//todo 处理失败..
		fmt.Println(err)
	} else {
		//todo 处理成功后通用操作
		pkt := g.ComResponse(packet.Id)
		connection.Write(pkt)
	}
}
