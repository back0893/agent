package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"github.com/back0893/goTcp/utils"
	"log"
)

/**
心跳不能被中控服务器说控制
因为这个是一个必要的服务
*/

type HeartBeatService struct {
}

func NewHeartBeatService() *HeartBeatService {
	return &HeartBeatService{}
}
func (m *HeartBeatService) Info() {
	pkt := src.NewPkt()
	pkt.Id = g.PING
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	if err := a.GetCon().Write(pkt); err != nil {
		log.Println(err)
	}
}
