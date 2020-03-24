package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

type CPUService struct {
	CurrentStatus string
	timerId       int64
}

func (m *CPUService) Cancel() {
	src.CancelTimer(m.timerId)
}

func (m *CPUService) GetCurrentStatus() string {
	return m.CurrentStatus
}

func (m *CPUService) SetCurrentStatus(status string) {
	m.CurrentStatus = status
}

func NewCPUService(status string) *CPUService {
	s := &CPUService{
		CurrentStatus: status,
	}
	s.upload(map[string]string{})
	return s
}

func (m *CPUService) Action(action string, args map[string]string) {
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

func (m *CPUService) Start(args map[string]string) error {
	m.CurrentStatus = "start"
	return nil
}

func (m *CPUService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m *CPUService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m CPUService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}
func (m *CPUService) upload(args map[string]string) {
	if m.timerId != 0 {
		src.CancelTimer(m.timerId)
	}
	interval := g.GetInterval(args, 30)

	m.timerId = src.AddTimer(interval*time.Second, func() {
		pkt := src.NewPkt()
		pkt.Id = g.CPU

		if m.Status(nil) == false {
			pkt.Data, _ = g.EncodeData("cpu service stop")
			return
		} else {
			if funcs.CpuPrepared() == false {
				funcs.UpdateCpuStat()
				return
			}

			err := funcs.UpdateCpuStat()
			if err != nil {
				log.Println(err)
				return
			}

			cpu := funcs.CpuMetrics()
			pkt.Data, err = g.EncodeData(cpu)
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
	})

}
func (m *CPUService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}
}
