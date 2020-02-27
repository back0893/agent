package main

import (
	"agent/src"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

type Server struct {
	*net.Server
	tw *src.TimingWheel
}

func NewServer() *Server {
	server := &Server{
		Server: net.NewServer(),
	}
	server.tw = src.NewTimingWheel(server.GetContext())
	return server
}
func (s *Server) Run() {
	go s.tw.Start()
	s.Server.Run()
}
func (s *Server) AddTimer(interval time.Duration, fn func()) int64 {
	return s.tw.AddTimer(time.Now(), interval, fn)
}
func (s *Server) TimerAt(when time.Time, fn func()) int64 {
	return s.tw.AddTimer(when, 0, fn)
}

func main() {
	utils.GlobalConfig.Load("json", "./app.json")
	server := NewServer()
	server.AddEvent(&src.Event{})
	server.AddProtocol(&src.Protocol{})
	server.AddTimer(2*time.Second, func() {
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
	server.Listen()
}
