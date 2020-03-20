package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"github.com/back0893/goTcp/utils"
	"log"
	"strconv"
	"strings"
)

type PortService struct {
	CurrentStatus string
	Ports         []int64
}

func (m *PortService) GetCurrentStatus() string {
	return m.CurrentStatus
}

func (m *PortService) SetCurrentStatus(status string) {
	m.CurrentStatus = status
}
func NewPortService() *PortService {
	return &PortService{
		Ports:         []int64{},
		CurrentStatus: "start",
	}
}
func (m *PortService) Action(action string, args map[string]string) {
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
	pkt.Data = []byte("启动memory")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m *PortService) Start(args map[string]string) error {
	m.CurrentStatus = "start"
	for _, port := range strings.Split(args["ports"], ",") {
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			continue
		}
		m.Ports = append(m.Ports, p)
	}
	return nil
}

func (m *PortService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m PortService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m PortService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}
func (m *PortService) upload() {
	pkt := src.NewPkt()
	pkt.Id = g.PortListen

	if m.Status(nil) == false {
		pkt.Data, _ = g.EncodeData("port service stop")
		return
	} else {
		ports, err := funcs.ListenTcpPortMetrics(m.Ports...)
		if err != nil {
			//todo 获得内存失败咋个处理
			log.Println(err)
			return
		}
		pkt.Data, err = g.EncodeData(ports)
		if err != nil {
			log.Println(err)
			return
		}
	}

	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
func (m *PortService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}
}
