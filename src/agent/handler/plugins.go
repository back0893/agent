package handler

import (
	"agent/src"
	"agent/src/agent/plugins"
	"agent/src/g/model"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
)

type Plugins struct {
}

func (p Plugins) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	decoder := gob.NewDecoder(bytes.NewReader(packet.Data))
	repPlugins := model.Plugins{}
	if err := decoder.Decode(&repPlugins); err != nil {
		log.Println(err)
		return
	}

	//因为git的拉取操作为一个耗时任务,
	//防止其他处理者被阻塞
	plugins.Git(utils.GlobalConfig.GetString("plugin.dir"), &repPlugins)
}
