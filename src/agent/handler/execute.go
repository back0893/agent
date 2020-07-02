package handler

import (
	"agent/src/agent/plugins"
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"context"
	"log"
	"path/filepath"
	"strings"

	iface2 "agent/src/agent/iface"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
)

type Execute struct {
}

func (e Execute) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	//执行某个特定的shell
	info := model.Execute{}
	var logID int32
	if err := g.DecodeData(packet.Data, &info, &logID); err != nil {
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
	plugins.PluginExecute(plugin, func(stdout, stderr *bytes.Buffer, err error, isTimeout bool) {
		var status int8 = 0
		var message string = ""
		if stderr != nil && stderr.String() != "" {
			message = string(stderr.Bytes())
		} else if isTimeout {
			// has be killed
			message = plugin.FilePath + "执行超时"
		} else if err != nil {
			message = err.Error()
		} else {
			// exec successfully
			message = string(stdout.Bytes())
			status = 1
		}
		pkt := g.ComResponse(logID, status, message)
		a := utils.GlobalConfig.Get(g.AGENT).(iface2.IAgent)
		if err := a.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	})
}
