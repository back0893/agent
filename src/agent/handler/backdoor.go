package handler

import (
	"agent/src/g"
	"bytes"
	"context"
	"log"
	"time"

	"github.com/back0893/goTcp/iface"
)

type BackDoor struct {
}

func (b BackDoor) Handler(ctx context.Context, packet *g.Packet, connection iface.IConnection) {
	//后台的执行shell命令
	//需要注意不要使用长时间等待的命令
	var shell = ""
	if err := g.DecodeData(packet.Data, &shell); err != nil {
		log.Println(err)
		return
	}
	cmd := g.Command{
		Name:    "bash",
		Args:    []string{"-c", shell},
		Timeout: 5 * 1000,
		Callback: func(stdout, stderr *bytes.Buffer, err error, isTimeout bool) {
			pkt := g.NewPkt()
			pkt.Id = g.BackDoor
			stderrStr := stderr.String()
			stdoutStr := stdout.String()
			if err != nil {
				pkt.Data, _ = g.EncodeData(err.Error())
			} else if isTimeout {
				pkt.Data, _ = g.EncodeData("执行超时")
			} else if stderrStr != "" {
				pkt.Data, _ = g.EncodeData(stderrStr)
			} else {
				pkt.Data, _ = g.EncodeData(stdoutStr)
			}
			connection.AsyncWrite(pkt, 5*time.Second)
		},
	}
	go cmd.Run()
}
