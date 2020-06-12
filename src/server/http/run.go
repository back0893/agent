package http

import (
	"agent/src/g"
	"encoding/json"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"io"
	"net/http"
	"time"
)

type BackDoor struct {
	Agent string `json:"agent"`
	Shell string `json:"shell"`
}

func WrapperRun(s iface.IServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		info := BackDoor{}
		body := json.NewDecoder(request.Body)
		if err := body.Decode(&info); err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		con, ok := g.GetCon(s, info.Agent)
		if ok == false {
			io.WriteString(writer, fmt.Sprintf("%s不存在或者没有上线", info.Agent))
			return
		}
		pkt := g.NewPkt()
		pkt.Id = g.BackDoor
		pkt.Data, _ = g.EncodeData(info.Shell)
		if err := con.AsyncWrite(pkt, 5*time.Second); err != nil {
			io.WriteString(writer, "发送失败")
			return
		}
	}
}
