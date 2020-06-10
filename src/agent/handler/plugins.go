package handler

import (
	"agent/src/agent/plugins"
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
)

type Plugins struct {
}

func (p Plugins) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	repPlugins := model.Plugins{}
	if err := g.DecodeData(packet.Data, &repPlugins); err != nil {
		log.Println(err)
		return
	}

	//因为git的拉取操作为一个耗时任务,
	//防止其他处理者被阻塞
	plugins.Git(utils.GlobalConfig.GetString("plugin.dir"), &repPlugins)
}
