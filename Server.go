package main

import (
	"agent/src"
	"agent/src/g"
	"agent/src/http"
	"agent/src/http/handler"
	"context"
	"flag"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
)

var (
	config string
)

func httpServer(ctx context.Context) {
	s := http.NewServer("0.0.0.0:9123")
	s.AddHandler("/", handler.SendTask)
	s.Run()
	select {
	case <-ctx.Done():
		s.Close(ctx)
	}
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
	go httpServer(server.GetContext())
	//todo http使用tcp连接上来,然后由这个转发给各个agent

	server.Listen(ip, port)
}
