package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
)

type HeartBeatService struct {
	CurrentStatus string //当前配置的服务状态
}

func (m *HeartBeatService) GetCurrentStatus() string {
	return m.CurrentStatus
}

func (m *HeartBeatService) SetCurrentStatus(status string) {
	m.CurrentStatus = status
}
func (m *HeartBeatService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}
	if m.Status(nil) == false {
		fmt.Printf("heart service stop")
		return
	}
	fmt.Println("ping")
	pkt := src.NewPkt()
	pkt.Id = g.PING
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	if err := a.GetCon().Write(pkt); err != nil {
		log.Println(err)
	}
}

func NewHeartBeatService() *HeartBeatService {
	return &HeartBeatService{
		CurrentStatus: "start",
	}
}

func (m *HeartBeatService) Action(action string, args map[string]string) {
	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
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
	pkt.Data = []byte("!启动心跳!")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m *HeartBeatService) Start(args map[string]string) error {
	m.CurrentStatus = "start"
	return nil
}

func (m *HeartBeatService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m *HeartBeatService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HeartBeatService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}
