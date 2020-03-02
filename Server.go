package main

import (
	"agent/src"
	"agent/src/g"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
)

var (
	config string
)

func mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
func init() {
	g.LoadInit()
}

func main() {
	flag.StringVar(&config, "c", "./app.json", "加载的配置json")
	flag.Parse()

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
