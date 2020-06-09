package iface

import (
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
)

type ServiceMethod interface {
	Handler(ctx context.Context, service *model.ServiceResponse, connection iface.IConnection) error
}
