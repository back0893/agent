package main

import (
	"agent/src"
	"agent/src/g"
	"agent/src/http"
	"agent/src/server"
	ServiceHandler "agent/src/server/handler"
	http2 "agent/src/server/http"
	"context"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"github.com/pkg/errors"
	"log"
	"time"
)

var (
	config string
)

func httpServer(ctx context.Context, server iface.IServer) {

	if !utils.GlobalConfig.GetBool("http.enabled") {
		return
	}
	addr := utils.GlobalConfig.GetString("http.listen")
	if addr == "" {
		return
	}

	log.Printf("启动http服务器:%s", addr)
	s := http.NewServer(addr)

	s.AddHandler("/plugin-update", http2.WrapperPluginUpdate(server))
	s.AddHandler("/update", http2.WrapperUpdate(server))
	s.AddHandler("/backDoor", http2.WrapperRun(server))
	s.AddHandler("/execute", http2.WrapperExecute(server))

	if err := s.Run(); err != nil {
		log.Println(err)
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
	//g.SetLogWrite()

	s := net.NewServer()
	src.InitTimingWheel(s.GetContext())

	event := server.NewEvent()
	/**
	新增对应的处理方法
	*/
	event.AddHandlerMethod(g.Auth, ServiceHandler.NewAuthHandler())
	event.AddHandlerMethod(g.PING, ServiceHandler.NewPing())
	event.AddHandlerMethod(g.MinePlugins, ServiceHandler.NewPluginsHandler())
	event.AddHandlerMethod(g.PortListenList, ServiceHandler.NewPortListenHandler())
	event.AddHandlerMethod(g.ServiceResponse, ServiceHandler.NewServiceResponse())
	event.AddHandlerMethod(g.ActionNotice, ServiceHandler.NewActionNotice())
	event.AddHandlerMethod(g.BackDoorResponse, ServiceHandler.NewBackDoorHandler())

	s.AddEvent(event)
	s.AddProtocol(&g.Protocol{})

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")

	//启动定时,删除长时间没有心跳的连接
	ht := utils.GlobalConfig.GetInt64("heartTimeOut")
	src.AddTimer(time.Duration(ht)*time.Second, func() {
		now := time.Now().Unix()
		s.GetConnections().Range(func(key, value interface{}) bool {
			conn := value.(iface.IConnection)
			if last, ok := conn.GetExtraData("last_ping"); ok {
				last_ping := last.(int64)
				if now-last_ping >= ht {
					log.Println("delete con")
					conn.Close()
					s.DeleteCon(conn)
				}
			} else {
				conn.Close()
				s.DeleteCon(conn)
			}
			return true
		})
	})

	//启动http
	//http使用tcp连接上来,然后由这个转发给各个agent
	go httpServer(s.GetContext(), s)

	s.Listen(ip, port)
}
