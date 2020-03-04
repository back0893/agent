package main

import (
	"agent/src"
	"agent/src/agent/model"
	"agent/src/g"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"io"
	"log"
	"net/http"
)

var (
	config string
)

func main() {
	flag.StringVar(&config, "c", "./app.json", "加载的配置json")
	flag.Parse()
	//加载
	g.LoadInit(config)

	server := net.NewServer()

	src.InitTimingWheel(server.GetContext())

	server.AddEvent(&src.Event{})
	server.AddProtocol(&src.Protocol{})

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")

	//启动一个channle http接受指令
	//y	actions:=make(chan *Action,20)
	taskQueue := src.NewTaskQueue()
	//启动http
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		action := src.Action{}
		body := json.NewDecoder(request.Body)
		if err := body.Decode(&action); err != nil {
			io.WriteString(writer, err.Error())
			return
		}
		taskQueue.Push(&action)
		writer.Write([]byte("ok"))
	})
	go http.ListenAndServe("0.0.0.0:9123", nil)

	//读取taskQueue,执行相应的操作
	go func() {
		for {
			action := taskQueue.Pop()
			fmt.Println(action.Action)
			cons := server.GetConnections()
			var icon iface.IConnection
			cons.Range(func(key, value interface{}) bool {
				con := value.(*net.Connection)
				v, ok := con.GetExtraData("auth")
				if ok == false {
					return true
				}
				auth := v.(*model.Auth)
				if auth.Username == action.DeviceId {
					icon = con
				}
				fmt.Println(action.DeviceId)
				fmt.Println(auth.Username)
				return false
			})
			if icon == nil {
				continue
			}
			var ipacket *src.Packet
			switch action.Action {
			case "start":
				fallthrough
			case "status":
				fallthrough
			case "stop":
				fallthrough
			case "restart":
				ipacket = src.NewPkt()
				ipacket.Id = g.Service
				b := bytes.NewBuffer([]byte{})
				encoder := gob.NewEncoder(b)
				service := model.Service{
					Service: "redis",
					Action:  action.Action,
				}
				encoder.Encode(service)
				ipacket.Data = b.Bytes()
			default:
				log.Println("新增失败,命令错误")
				continue
			}
			icon.Write(ipacket)
		}
	}()

	server.Listen(ip, port)
}
