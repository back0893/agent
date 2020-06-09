package server

import (
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

func NewEvent() *Event {
	return &Event{
		g.NewEvent(),
	}
}

type Event struct {
	*g.Event
}

func (e *Event) OnConnect(ctx context.Context, connection iface.IConnection) {
	log.Printf("连接%d", connection.GetId())
}

func (Event) OnClose(ctx context.Context, connection iface.IConnection) {
	if v, ok := connection.GetExtraData("auth"); ok {
		auth := v.(*model.Auth)
		log.Printf("用户%s断开连接", auth.Username)
	}
}
