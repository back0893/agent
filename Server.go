package main

import (
	"agent/src"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"time"
)

func TimerFn(server *net.Server) {
	heart := utils.GlobalConfig.GetInt("heartInterval")
	timeOut := utils.GlobalConfig.GetInt64("heartTimeOut")
	ticker := time.NewTicker(time.Second * time.Duration(heart))
	for {
		select {
		case now := <-ticker.C:
			conMap := server.GetConnections()
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
					connection.Close()
				}
				return true
			})
		}
	}
}

func main() {
	utils.GlobalConfig.Load("json", "./app.json")
	server := net.NewServer()
	server.AddEvent(&src.Event{})
	server.AddProtocol(&src.Protocol{})
	go TimerFn(server)
	server.Listen()
}
