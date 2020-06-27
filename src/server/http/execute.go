package http

import (
	"agent/src/g"
	"agent/src/g/model"
	"encoding/json"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"io"
	"net/http"
	"time"
)

type execute struct {
	TimeOut int    `json:"timeout"`
	Shell   string `json:"shell"`
	Agent   string `json:"agent"`
}

func WrapperExecute(s iface.IServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		info := execute{}
		body := json.NewDecoder(request.Body)
		if err := body.Decode(&info); err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		if info.TimeOut < 1 {
			info.TimeOut = 1
		}
		con, ok := g.GetCon(s, info.Agent)
		if ok == false {
			io.WriteString(writer, fmt.Sprintf("%s不存在或者没有上线", info.Agent))
			return
		}
		pkt := g.NewPkt()
		pkt.Id = g.Execute
		ex := model.Execute{
			File:    info.Shell,
			TimeOut: info.TimeOut,
		}
		pkt.Data, _ = g.EncodeData(ex)
		if err := con.AsyncWrite(pkt, 5*time.Second); err != nil {
			io.WriteString(writer, "发送失败")
			return
		}
	}
}
