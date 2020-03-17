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

type LoadAvgServiceService struct {
	CurrentStatus string
}

func NewLoadAvgServiceService() *LoadAvgServiceService {
	return &LoadAvgServiceService{
		CurrentStatus: "start",
	}
}
func (m *LoadAvgServiceService) Action(action string, args map[string]string) {
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
	pkt.Data = []byte("启动负载")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m *LoadAvgServiceService) Start(args map[string]string) error {
	m.CurrentStatus = "start"

	return nil
}

func (m *LoadAvgServiceService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m *LoadAvgServiceService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m LoadAvgServiceService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}

func (m *LoadAvgServiceService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}
	if m.Status(nil) == false {
		fmt.Sprintf("loadAvg  service stop")
		return
	}
	loadAvg, err := funcs.LoadAvgMetrics()
	if err != nil {
		//todo 获得内存失败咋个处理
		log.Println(err)
		return
	}
	pkt := src.NewPkt()
	pkt.Id = g.LoadAvg
	pkt.Data, err = g.EncodeData(loadAvg)
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
