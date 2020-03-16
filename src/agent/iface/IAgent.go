package iface

import (
	"context"
	"github.com/back0893/goTcp/iface"
	"sync"
)

type IAgent interface {
	GetCon() iface.IConnection
	GetWaitGroup() *sync.WaitGroup
	GetContext() context.Context
	AddEvent(event iface.IEvent)
	AddConnect(fn func(context.Context, iface.IConnection))
	AddClose(fn func(context.Context, iface.IConnection))
	AddProtocol(protocol iface.IProtocol)
	Start()
	IsStop() bool
	Stop()
	Wait()
	ReCon(ctx context.Context, con iface.IConnection)
}
