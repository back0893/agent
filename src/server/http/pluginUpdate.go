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

type action struct {
	Agent string   `json:"agent"` //发送给agent的名称
	Git   []string `json:"git"`   //拉取的git地址
}

/**
为了使http处理人获得net.server的参数,有不想用全局变量
*/
func WrapperPluginUpdate(s iface.IServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		action := &action{}
		body := json.NewDecoder(request.Body)
		if err := body.Decode(&action); err != nil {
			io.WriteString(writer, "json错误:"+err.Error())
			return
		}
		con, ok := g.GetCon(s, action.Agent)
		if ok == false {
			io.WriteString(writer, fmt.Sprintf("%s不存在或者没有上线", action.Agent))
			return
		}
		pkt := g.NewPkt()
		pkt.Id = g.MinePluginsResponse
		plugins := model.Plugins{
			Uri: action.Git,
		}
		pkt.Data, _ = g.EncodeData(plugins)
		if err := con.AsyncWrite(pkt, time.Second*2); err != nil {
			io.WriteString(writer, "发送至agent超时")
			return
		}
		//todo 连接到tcp发送消息
		io.WriteString(writer, "success")
	}
}
