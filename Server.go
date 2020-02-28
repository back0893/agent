package main

import (
	"agent/src"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

func main() {
	utils.GlobalConfig.Load("json", "./app.json")
	server := net.NewServer()
	src.InitTimingWheel(server.GetContext())

	server.AddEvent(&src.Event{})
	server.AddProtocol(&src.Protocol{})
	src.AddTimer(2*time.Second, func() {
		conMap := server.GetConnections()
		timeOut := utils.GlobalConfig.GetInt64("heartTimeOut")
		now := time.Now()
		nowTimestamp := now.Unix()
		conMap.Range(func(key, value interface{}) bool {
			connection := value.(*net.Connection)
			t, ok := connection.GetExtraData("heart")
			if ok == false {
				//可能是刚连接上,还未发送消息
				connection.SetExtraData("heart", now.Unix())
				return true
			}
			headerTimestamp := t.(int64)
			if nowTimestamp-headerTimestamp > timeOut {
				//当前连接的心跳灭有正常上报.判断是网络连接出现问题
				log.Println("心跳过期,主动关闭", connection.GetId())
				connection.Close()
			}
			return true
		})
	})

	//主动连接到一个任务发送系统,等待任务下达
	//任务下达完成后将通知任务系统任务结果
	//TODO task_client

	ip := utils.GlobalConfig.GetString("Ip")
	port := utils.GlobalConfig.GetInt("Port")
	server.Listen(ip, port)
}
