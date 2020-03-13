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
var cpuId int64

type CPUService struct {
	agent iface.IAgent
}

func NewCPUService(agent iface.IAgent) *CPUService {
	return &CPUService{agent: agent}
}

func (m *CPUService) Action(action string, args []string) {
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

func (m CPUService) Start(args []string) error {
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	cpuId = src.AddTimer(time.Duration(num)*time.Second, func() {

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

		pkt := src.NewPkt()
		pkt.Id = g.CPU
		pkt.Data, err = g.EncodeData(cpu)
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

func (m CPUService) Stop([]string) error {
	if cpuId > 0 {
		src.CancelTimer(cpuId)
	}
	return nil
}

func (m CPUService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m CPUService) Status([]string) bool {
	if cpuId > 0 {
		return true
	}
	return false
}
