package g

import (
	"context"
	"github.com/back0893/goTcp/iface"
)

type IHandlerMethod interface {
	Handler(ctx context.Context, packet *Packet, connection iface.IConnection)
}
