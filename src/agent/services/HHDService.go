package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

type HHDService struct {
	CurrentStatus int
	timeId        int64
}

func NewHHDService(status int) *HHDService {
	s := &HHDService{
		CurrentStatus: status,
	}
	s.Upload(map[string]string{})
	return s
}
func (m *HHDService) GetCurrentStatus() int {
	return m.CurrentStatus
}

func (m *HHDService) SetCurrentStatus(status int) {
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
	if m.Status(nil) {
		return nil
	}
	m.CurrentStatus = 1
	return nil
}

func (m *HHDService) Stop(map[string]string) error {
	m.CurrentStatus = 0
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
	return m.CurrentStatus == 1
}
func (m HHDService) info() {
	info := model.NewServiceResponse(g.HHD, m.CurrentStatus)
	if m.Status(nil) == false {
		info.Status = 0
		info.Info = "启动失败"
	} else {
		disks, err := funcs.DiskUseMetrics()
		if err != nil {
			//todo 获得失败咋个处理
			log.Println(err)
			return
		}
		info.Info = disks
	}
	pkt := src.ServiceResponsePkt(info)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
func (m *HHDService) Upload(args map[string]string) {
	if m.timeId != 0 {
		src.CancelTimer(m.timeId)
	}
	interval := g.GetInterval(args, 10)
	m.timeId = src.AddTimer(interval*time.Second, m.info)
}
func (m *HHDService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == 0 {
		m.CurrentStatus = 1
	} else if m.CurrentStatus == 1 && run == false {
		m.Start(map[string]string{})
	}
}
func (m *HHDService) Cancel() {
	src.CancelTimer(m.timeId)
}
