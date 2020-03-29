package iface

import (
	"agent/src"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
)

type HandlerMethod interface {
	Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection)
}

type ServiceMethod interface {
	Handler(ctx context.Context, service *model.ServiceResponse, connection iface.IConnection) error
}
