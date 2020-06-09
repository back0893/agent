package g

import (
	"agent/src"
	iface2 "agent/src/g/iface"
	"agent/src/server/handler"
	"context"
	"github.com/back0893/goTcp/iface"
	"sync"
)

func NewEvent() *Event {
	e := &Event{
		methods: make(map[int32]iface2.IHandlerMethod),
	}
	e.AddHandlerMethod(0, &handler.DefaultMethod{})
	return e
}

type Event struct {
	lock    sync.RWMutex
	methods map[int32]iface2.IHandlerMethod
}

func (e *Event) AddHandlerMethod(id int32, fn iface2.IHandlerMethod) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.methods[id] = fn
}
func (e *Event) GetMethod(id int32) iface2.IHandlerMethod {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fn, ok := e.methods[id]
	if ok {
		return fn
	}
	return e.methods[0]
}
func (e *Event) OnConnect(context.Context, iface.IConnection) {
}

func (e *Event) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	pkt := packet.(*src.Packet)
	id := pkt.Id
	fn := e.GetMethod(id)
	fn.Handler(ctx, pkt, connection)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
}