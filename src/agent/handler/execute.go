package handler

import (
	"agent/src/agent/plugins"
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
)

type Execute struct {
}

func (e Execute) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	//执行某个特定的shell
	info := model.Execute{}
	if err := g.DecodeData(packet.Data, &info); err != nil {
		log.Println(err)
		return
	}

	//如果不为/ 开始说明是相对路径
	if strings.Index(info.File, "/") != 0 {
		info.File = filepath.Join(utils.GlobalConfig.GetString("plugin.dir"), info.File)
	}
	plugin := &plugins.Plugin{
		FilePath: info.File,
		Interval: info.TimeOut,
		IsRepeat: false,
		MTime:    0,
	}
	plugins.PluginRun(plugin)
	pkt := g.ComResponse(packet.Id)
	connection.AsyncWrite(pkt, 5*time.Second)
}
