package main

import (
	"agent/src"
	"agent/src/g"
	"agent/src/http"
	"agent/src/http/handler"
	"agent/src/server"
	"context"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"log"
)

var (
	config string
)

func httpServer(ctx context.Context, server iface.IServer) {
	host := utils.GlobalConfig.GetString("http.host")
	port := utils.GlobalConfig.GetInt("http.port")
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("启动http服务器:%s", addr)
	s := http.NewServer(addr)

	s.AddHandler("/sendTask", handler.WrapperSendTask(server))
	s.AddHandler("/update", handler.WrapperUpdate(server))

	go func() {
		if err := s.Run(); err != nil {
			panic(err)
		}
	}()
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

	s := net.NewServer()
	src.InitTimingWheel(s.GetContext())

	event := &server.Event{}
	s.AddEvent(event)
	s.AddProtocol(&src.Protocol{})

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")

	//启动http
	//http使用tcp连接上来,然后由这个转发给各个agent
	go httpServer(s.GetContext(), s)

	s.Listen(ip, port)
}
