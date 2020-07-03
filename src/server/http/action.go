package http

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/back0893/goTcp/iface"
)

type Action struct {
	Code    string `json:"code"`
	ID      string `json:"id"`
	Command string `json:"command"`
	LogID   int32  `json:"logId"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	server, ok := r.Context().Value(g.SERVER).(iface.IServer)
	mm := r.Context().Value("test").(string)
	log.Println(mm)

	if !ok {
		w.Write([]byte("!错误!"))
		return
	}
	info := Action{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&info); err != nil {
		w.Write([]byte("!解析json错误!"))
		return
	}

	pkt := g.NewPkt()
	switch info.Code {
	case "0001":
		//重启
		pkt.Id = g.STOP
	case "0003":
		//代码更新
		pkt.Id = g.MinePlugins
		data := model.Plugins{
			Uri: []string{info.Command},
		}
		pkt.Data, _ = g.EncodeData(data)
	case "0004":
		//监听端口
		pkt.Id = g.PortListenList
		data := make([]int32, 0)
		portStr := strings.Split(info.Command, ",")

		for _, p := range portStr {
			if port, err := strconv.Atoi(p); err != nil {
				continue
			} else {
				data = append(data, int32(port))
			}
		}

		pkt.Data, _ = g.EncodeData(data)
	case "0005":
		//更新agent
		pkt.Id = g.UPDATE
		data := model.UpdateInfo{
			URL:  info.Command,
			Type: 1,
		}
		pkt.Data, _ = g.EncodeData(data)
	case "0006":
		//更新配置
		pkt.Id = g.UPDATE
		data := model.UpdateInfo{
			URL:  info.Command,
			Type: 2,
		}
		pkt.Data, _ = g.EncodeData(data)
	case "0007":
		//执行一个脚本
		pkt.Id = g.Execute
		data := model.Execute{
			File:    info.Command,
			TimeOut: 10 * 1000,
		}
		pkt.Data, _ = g.EncodeData(data)
	default:
		w.Write([]byte("未知执行,请确认"))
		return
	}

	con, ok := g.GetCon(server, info.ID)
	if ok == false {
		ep, ok := db.DbConnections.Get("ep")
		if !ok {
			io.WriteString(w, "数据库连接失败!")
		}
		ep.Exec("update cc_service_log set status=?,log=? where id=?", 3, "不存在或者没有上线", info.LogID)
		io.WriteString(w, fmt.Sprintf("%s不存在或者没有上线", info.ID))
		return
	}
	//新增logId
	logID, err := g.EncodeData(info.LogID)
	if err != nil {
		w.Write([]byte("非法logId"))
		return
	}
	pkt.Data = append(pkt.Data, logID...)
	if err := con.AsyncWrite(pkt, time.Second*5); err != nil {
		ep, ok := db.DbConnections.Get("ep")
		if !ok {
			io.WriteString(w, "数据库连接失败!")
		}
		ep.Exec("update cc_service_log set status=?,log=? where id=?", 3, "发送超时", info.LogID)
		io.WriteString(w, "发送超时")
		return
	}
	//todo 连接到tcp发送消息
	io.WriteString(w, "success")

}
