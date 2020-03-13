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
var loadId int64

type LoadAvgServiceService struct {
	agent iface.IAgent
}

func NewLoadAvgServiceService(agent iface.IAgent) *HHDService {
	return &HHDService{agent: agent}
}
func (m *LoadAvgServiceService) Action(action string, args []string) {
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

func (m LoadAvgServiceService) Start(args []string) error {
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	loadId = src.AddTimer(time.Duration(num)*time.Second, func() {
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
		err = m.agent.GetCon().Write(pkt)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return nil
}

func (m LoadAvgServiceService) Stop([]string) error {
	if loadId > 0 {
		src.CancelTimer(loadId)
	}
	return nil
}

func (m LoadAvgServiceService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m LoadAvgServiceService) Status([]string) bool {
	if loadId > 0 {
		return true
	}
	return false
}
