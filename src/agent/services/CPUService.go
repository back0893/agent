package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"errors"
	"github.com/back0893/goTcp/utils"
	"log"
	"strconv"
	"time"
)

//一个私有的全局变量
var cpuId int64

type CPUService struct {
}

func NewCPUService() *CPUService {
	return &CPUService{}
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

func (m CPUService) Start(args map[string]string) error {
	if m.Status(args) {
		return errors.New("service已经启动")
	}
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args["interval"])
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
		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		err = a.GetCon().Write(pkt)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return nil
}

func (m CPUService) Stop(map[string]string) error {
	if cpuId > 0 {
		src.CancelTimer(cpuId)
	}
	return nil
}

func (m CPUService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m CPUService) Status(map[string]string) bool {
	if cpuId > 0 {
		return true
	}
	return false
}
