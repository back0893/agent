package handler

import (
	"agent/src"
	"agent/src/agent/funcs"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/back0893/goTcp/iface"
	"log"
)

type Ports struct {
}

func (p Ports) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	decoder := gob.NewDecoder(bytes.NewReader(packet.Data))
	ports := make([]int64, 0)
	if err := decoder.Decode(&ports); err != nil {
		log.Println(err)
	}
	funcs.AppendPorts(ports...)
}
