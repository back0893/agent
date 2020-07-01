package http

import (
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"encoding/json"
	"log"
	"net/http"
)

/**
中转,脚本也可以请求http接口发送数据,已取代标准输出
*/
func Hanlder(w http.ResponseWriter, req *http.Request) {
	agent := req.Context().Value(g.SERVER).(iface.IAgent)
	if req.ContentLength == 0 {
		http.Error(w, "body is blank", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var metrics []*model.MetricValue
	err := decoder.Decode(&metrics)
	if err != nil {
		http.Error(w, "connot decode body,err:"+err.Error(), http.StatusBadRequest)
		return
	}
	pkt := g.NewPkt()
	pkt.Id = g.ActionNotice
	if pkt.Data, err = g.EncodeData(metrics); err != nil {
		log.Print(err)
		return
	}
	if err := agent.GetCon().Write(pkt); err != nil {
		log.Println(err)
	}
	w.Write([]byte("success"))
}
