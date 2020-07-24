package handler

import (
	iface2 "agent/src/agent/iface"
	"agent/src/g"
	"context"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
)

type Stop struct{}

//Handler 停止处理
func (s Stop) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	agent := utils.GlobalConfig.Get(g.AGENT).(iface2.IAgent)
	agent.Stop()
}
