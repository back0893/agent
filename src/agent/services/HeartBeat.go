package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"fmt"
	"log"
	"strconv"
	"time"
)

//一个私有的全局变量
var heartId int64

type HeartBeatService struct {
	agent iface.IAgent
}

func NewHeartBeatService(agent iface.IAgent) *HeartBeatService {
	return &HeartBeatService{agent: agent}
}

func (m *HeartBeatService) Action(action string, args []string) {
	switch action {
	case "start":
		m.Start(args)
	case "stop":
		m.Stop(args)
	case "restart":
		m.Restart(args)
	case "status":
		m.Status(args)
	}
	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	pkt.Data = []byte("!启动心跳!")
	err := m.agent.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m HeartBeatService) Start(args []string) error {
	var num = 10
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	fmt.Println(num)
	heartId = src.AddTimer(time.Duration(num)*time.Second, func() {
		fmt.Println("ping")
		pkt := src.NewPkt()
		pkt.Id = g.PING
		if err := m.agent.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	})
	return nil
}

func (m HeartBeatService) Stop([]string) error {
	if heartId > 0 {
		src.CancelTimer(heartId)
	}
	return nil
}

func (m HeartBeatService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HeartBeatService) Status([]string) bool {
	if heartId > 0 {
		return true
	}
	return false
}
