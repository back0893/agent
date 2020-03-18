package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
)

type HHDService struct {
	CurrentStatus string
}

func NewHHDService() *HHDService {
	return &HHDService{
		CurrentStatus: "start",
	}
}
func (m *HHDService) GetCurrentStatus() string {
	return m.CurrentStatus
}

func (m *HHDService) SetCurrentStatus(status string) {
	m.CurrentStatus = status
}
func (m *HHDService) Action(action string, args map[string]string) {
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
	pkt.Data = []byte("启动硬盘")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m *HHDService) Start(args map[string]string) error {
	m.CurrentStatus = "start"
	return nil
}

func (m *HHDService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m HHDService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HHDService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}
func (m *HHDService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}

	if m.Status(nil) == false {
		fmt.Printf("hhd service stop")
		return
	}

	disks, err := funcs.DiskUseMetrics()
	if err != nil {
		//todo 获得内存失败咋个处理
		log.Println(err)
		return
	}
	pkt := src.NewPkt()
	pkt.Id = g.HHD
	pkt.Data, err = g.EncodeData(disks)
	if err != nil {
		log.Println(err)
		return
	}
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err = a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
