package handler

import (
	"agent/src/g"
	model2 "agent/src/g/model"
	"encoding/json"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"io"
	"net/http"
)

type action struct {
	Agent   string            //发送给agent的名称
	Service int               //发送的服务名称
	Action  string            //服务对应的动作
	Args    map[string]string //对应的传递参数
}

/**
为了使http处理人获得net.server的参数,有不想用全局变量
*/
func WrapperSendTask(s iface.IServer) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		action := &action{}
		body := json.NewDecoder(request.Body)
		if err := body.Decode(&action); err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		con, ok := g.GetCon(s, action.Agent)
		if ok == false {
			io.WriteString(writer, fmt.Sprintf("%s不存在或者没有上线", action.Agent))
			return
		}
		pkt := g.ServicePkt(model2.NewService(action.Service, action.Action, action.Args))
		con.Write(pkt)
		//todo 连接到tcp发送消息
		writer.Write([]byte("ok"))
	}
}
