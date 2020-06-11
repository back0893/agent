package g

import (
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
	"sync"
)

type DefaultMethod struct {
}

func NewDefaultMethod() *DefaultMethod {
	return &DefaultMethod{}
}

func (d DefaultMethod) Handler(ctx context.Context, packet *Packet, connection iface.IConnection) {
	log.Printf("方法还未被时实现method_id===>%d", packet.Id)
}

func NewEvent() *Event {
	e := &Event{
		methods: make(map[int32]IHandlerMethod),
	}
	e.AddHandlerMethod(0, &DefaultMethod{})
	return e
}

type Event struct {
	lock    sync.RWMutex
	methods map[int32]IHandlerMethod
}

func (e *Event) AddHandlerMethod(id int32, fn IHandlerMethod) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.methods[id] = fn
}
func (e *Event) GetMethod(id int32) IHandlerMethod {
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
	pkt := packet.(*Packet)
	id := pkt.Id
	fn := e.GetMethod(id)
	log.Println(id)
	fn.Handler(ctx, pkt, connection)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
}
