package handler

import (
	"agent/src/agent/plugins"
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"log"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
)

type Plugins struct {
}

func (p Plugins) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	repPlugins := model.Plugins{}
	var logID int32
	if err := g.DecodeData(packet.Data, &repPlugins, &logID); err != nil {
		log.Println(err)
		return
	}

	//因为git的拉取操作为一个耗时任务,
	//防止其他处理者被阻塞
	plugins.Git(utils.GlobalConfig.GetString("plugin.dir"), &repPlugins, logID)

}
