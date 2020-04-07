package server

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/handler"
	serverFace "agent/src/server/iface"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
	"sync"
)

func NewEvent() *Event {
	e := &Event{
		methods: make(map[int32]serverFace.HandlerMethod),
	}
	e.AddHandlerMethod(0, &handler.DefaultMethod{})
	return e
}

type Event struct {
	lock    sync.RWMutex
	methods map[int32]serverFace.HandlerMethod
}

func (e *Event) AddHandlerMethod(id int32, fn serverFace.HandlerMethod) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	e.methods[id] = fn
}
func (e *Event) GetMethod(id int32) serverFace.HandlerMethod {
	e.lock.RLock()
	defer e.lock.RUnlock()
	fn, ok := e.methods[id]
	if ok {
		return fn
	}
	return e.methods[0]
}
func (e *Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	SetTimeOut(connection.GetRawCon())
}

func (e *Event) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	SetTimeOut(connection.GetRawCon())
	pkt := packet.(*src.Packet)
	id := pkt.Id
	//如果未认证过,断开
	if id != g.Auth {
		if _, ok := connection.GetExtraData("auth"); !ok {
			connection.Close()
			return
		}
	}
	fn := e.GetMethod(id)
	fn.Handler(ctx, pkt, connection)
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
	if v, ok := connection.GetExtraData("auth"); ok {
		auth := v.(*model.Auth)
		log.Printf("用户%s断开连接", auth.Username)
	}
}
