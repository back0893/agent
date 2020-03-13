package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"log"
	"strconv"
	"time"
)

//一个私有的全局变量
var memId int64

type MemoryService struct {
	agent iface.IAgent
}

func NewMemoryService(agent iface.IAgent) *MemoryService {
	return &MemoryService{agent: agent}
}

func (m *MemoryService) Action(action string, args []string) {
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
	err := m.agent.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m MemoryService) Start(args []string) error {
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	memId = src.AddTimer(time.Duration(num)*time.Second, func() {
		memory, err := funcs.MemMetrics()
		if err != nil {
			//todo 获得内存失败咋个处理
			log.Println(err)
			return
		}
		pkt := src.NewPkt()
		pkt.Id = g.MEM
		pkt.Data, err = g.EncodeData(memory)
		if err != nil {
			log.Println(err)
			return
		}
		err = m.agent.GetCon().Write(pkt)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return nil
}

func (m MemoryService) Stop([]string) error {
	if memId > 0 {
		src.CancelTimer(memId)
	}
	return nil
}

func (m MemoryService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m MemoryService) Status([]string) bool {
	if memId > 0 {
		return true
	}
	return false
}
