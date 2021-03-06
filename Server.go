package main

import (
	"agent/src"
	"agent/src/g"
	"agent/src/http"
	"agent/src/http/handler"
	"agent/src/server"
	ServiceHandler "agent/src/server/handler"
	"context"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"github.com/pkg/errors"
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
			log.Println(err)
		}
	}()
	select {
	case <-ctx.Done():
		s.Close(ctx)
	}
}
func main() {
	defer func() {
		if err := recover(); err != nil {
			stackerr := errors.WithStack(errors.New(fmt.Sprintln(err)))
			log.Printf("%+v", stackerr)
		}
	}()
	flag.StringVar(&config, "c", "./app.json", "加载的配置json")
	flag.Parse()
	//加载
	g.LoadInit(config)

	s := net.NewServer()
	src.InitTimingWheel(s.GetContext())

	event := server.NewEvent()
	/**
	新增对应的处理方法
	*/
	event.AddHandlerMethod(g.Auth, ServiceHandler.NewAuthHandler())
	event.AddHandlerMethod(g.ServiceResponse, ServiceHandler.NewServiceResponse())
	event.AddHandlerMethod(g.PING, ServiceHandler.NewPing())
	event.AddHandlerMethod(0, ServiceHandler.NewDefaultMethod())

	s.AddEvent(event)
	s.AddProtocol(&src.Protocol{})

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")

	//启动http
	//http使用tcp连接上来,然后由这个转发给各个agent
	go httpServer(s.GetContext(), s)

	s.Listen(ip, port)
}
