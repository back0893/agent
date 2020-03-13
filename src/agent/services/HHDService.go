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
var hhdId int64

type HHDService struct {
	agent iface.IAgent
}

func NewHHDService(agent iface.IAgent) *HHDService {
	return &HHDService{agent: agent}
}

func (m *HHDService) Action(action string, args []string) {
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

func (m HHDService) Start(args []string) error {
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	hhdId = src.AddTimer(time.Duration(num)*time.Second, func() {
		disks, err := funcs.DiskUseMetrics()
		if err != nil {
			//todo 获得内存失败咋个处理
			log.Println(err)
			return
		}
		pkt := src.NewPkt()
		pkt.Id = g.HHD
		pkt.Data, err = g.EncodeData(disks)
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

func (m HHDService) Stop([]string) error {
	if hhdId > 0 {
		src.CancelTimer(hhdId)
	}
	return nil
}

func (m HHDService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HHDService) Status([]string) bool {
	if hhdId > 0 {
		return true
	}
	return false
}
