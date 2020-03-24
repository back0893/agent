package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

type MemoryService struct {
	CurrentStatus string
	timeId        int64
}

func (m *MemoryService) GetCurrentStatus() string {
	return m.CurrentStatus
}

func (m *MemoryService) SetCurrentStatus(status string) {
	m.CurrentStatus = status
}
func NewMemoryService(status string) *MemoryService {
	s := &MemoryService{
		CurrentStatus: status,
	}
	s.upload(map[string]string{})
	return s
}

func (m *MemoryService) Action(action string, args map[string]string) {
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

func (m *MemoryService) Start(args map[string]string) error {
	m.CurrentStatus = "start"
	return nil
}

func (m *MemoryService) Stop(map[string]string) error {
	m.CurrentStatus = "stop"
	return nil
}

func (m MemoryService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m MemoryService) Status(map[string]string) bool {
	return m.CurrentStatus == "start"
}
func (m *MemoryService) upload(args map[string]string) {
	if m.timeId != 0 {
		src.CancelTimer(m.timeId)
	}

	m.timeId = src.AddTimer(g.GetInterval(args, 30)*time.Second, func() {
		pkt := src.NewPkt()
		pkt.Id = g.MEM
		if m.Status(nil) == false {
			pkt.Data, _ = g.EncodeData("memeroy service stop")
		} else {
			memory, err := funcs.MemMetrics()
			if err != nil {
				//todo 获得内存失败咋个处理
				log.Println(err)
				return
			}
			pkt.Data, err = g.EncodeData(memory)
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
func (m *MemoryService) Watcher() {
	run := m.Status(nil)
	if run == true && m.CurrentStatus == "end" {
		m.CurrentStatus = "start"
	} else if m.CurrentStatus == "start" && run == false {
		m.Start(map[string]string{})
	}
}
func (m *MemoryService) Cancel() {
	fmt.Println("!!mem cancel!!")
	src.CancelTimer(m.timeId)
}
