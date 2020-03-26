package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"github.com/back0893/goTcp/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type PortService struct {
	CurrentStatus int
	Ports         []int64
	timeId        int64
}

func (m *PortService) GetCurrentStatus() int {
	return m.CurrentStatus
}

func (m *PortService) SetCurrentStatus(status int) {
	m.CurrentStatus = status
}
func NewPortService(status int) *PortService {
	s := &PortService{
		Ports:         []int64{},
		CurrentStatus: status,
	}
	s.Upload(map[string]string{})
	return s
}
func (m *PortService) Action(action string, args map[string]string) {
	switch action {
	case "start":
		m.Start(args)
	case "stop":
		m.Stop(args)
	case "restart":
		m.Restart(args)
	case "listen":
		m.Ports = m.SplitPort(args)
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
func (m *PortService) SplitPort(args map[string]string) []int64 {
	portSlice := make([]int64, 0)
	ports, ok := args["ports"]
	if ok == false {
		return portSlice
	}
	for _, port := range strings.Split(ports, ",") {
		p, err := strconv.ParseInt(port, 10, 64)
		if err != nil {
			continue
		}
		portSlice = append(portSlice, p)
	}
	return portSlice
}
func (m *PortService) Start(args map[string]string) error {
	if m.Status(nil) {
		return nil
	}
	m.CurrentStatus = 1
	m.Ports = m.SplitPort(args)
	return nil
}

func (m *PortService) Stop(map[string]string) error {
	m.CurrentStatus = 0
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
	return m.CurrentStatus == 1
}
func (m PortService) info(args map[string]string) {
	info := model.NewServiceResponse(g.PortListen, m.CurrentStatus)
	portListen := m.Ports
	if args != nil {
		portListen = m.SplitPort(args)
	}
	if len(portListen) == 0 {
		return
	}
	result, err := funcs.ListenTcpPortMetrics(portListen...)
	if err != nil {
		//todo 获得失败咋个处理
		log.Println(err)
		return
	}
	info.Info, err = g.EncodeData(result)
	pkt := src.ServiceResponsePkt(info)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err = a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
func (m *PortService) Upload(args map[string]string) {
	if m.timeId != 0 {
		src.CancelTimer(m.timeId)
	}
	m.timeId = src.AddTimer(g.GetInterval(args, 30)*time.Second, func() {
		m.info(nil)
	})
}
func (m *PortService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == 0 {
		m.CurrentStatus = 1
	} else if m.CurrentStatus == 1 && run == false {
		m.Start(map[string]string{})
	}
}
func (m *PortService) Cancel() {
	src.CancelTimer(m.timeId)
}
