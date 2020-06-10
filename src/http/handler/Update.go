package handler

import (
	"agent/src/g"
	"encoding/json"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"io"
	"net/http"
)

type UpdateInfo struct {
	Version int    //版本号
	Url     string //更新地址
	Agent   string
}

func WrapperUpdate(s iface.IServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		info := &UpdateInfo{}
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
		pkt.Id = g.UPDATE
		pkt.Data, _ = g.EncodeData(info)
		con.Write(pkt)
		writer.Write([]byte("ok"))
	}
}
