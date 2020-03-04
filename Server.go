package main

import (
	"agent/src"
	"agent/src/g"
	"agent/src/http"
	"agent/src/http/handler"
	"flag"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
)

var (
	config string
)

func httpServer() {
	s := http.NewServer()
	s.AddHandler("/", handler.SendTask)
	s.Run()
}
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

	//启动http
	go httpServer()

	server.Listen(ip, port)
}
