package main

import (
	"agent/src"
	"fmt"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"github.com/spf13/cast"
	"log"
	"os"
	"strings"
	"time"
)

func mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
func init() {
	var path string
	p := utils.GlobalConfig.Get("runtime")
	if p == nil {
		path = "./runtime"
	} else {
		path = cast.ToString(p)
	}
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) == false {
			err := mkdir(path)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	filePath := fmt.Sprintf("%s/%s.log", strings.TrimRight(path, "/"), time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)
}
func main() {
	utils.GlobalConfig.Load("json", "./app.json")
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
