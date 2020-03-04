package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

type Action struct {
	agent   string //发送给agent的名称
	service string //发送的服务名称
	action  string //服务对应的动作
}

func SendTask(writer http.ResponseWriter, request *http.Request) {
	action := &Action{}
	body := json.NewDecoder(request.Body)
	if err := body.Decode(&action); err != nil {
		io.WriteString(writer, err.Error())
		return
	}
	//todo 连接到tcp发送消息
	writer.Write([]byte("ok"))
}
