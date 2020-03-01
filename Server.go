package main

import (
	"agent/src"
	"agent/src/g"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
)

func init() {
	g.LoadInit()
}

func main() {
	server := net.NewServer()
	src.InitTimingWheel(server.GetContext())

	server.AddEvent(&src.Event{})
	server.AddProtocol(&src.Protocol{})

	//主动连接到一个任务发送系统,等待任务下达
	//任务下达完成后将通知任务系统任务结果
	//TODO task_client

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")
	server.Listen(ip, port)
}
