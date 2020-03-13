package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"log"
	"strconv"
	"strings"
	"time"
)

//一个私有的全局变量
var portId int64

type PortService struct {
	agent iface.IAgent
}

func NewPortService(agent iface.IAgent) *PortService {
	return &PortService{agent: agent}
}
func (m *PortService) Action(action string, args []string) {
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

func (m PortService) Start(args []string) error {
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args[0])
		if err == nil {
			num = n
		}
	}
	portId = src.AddTimer(time.Duration(num)*time.Second, func() {
		//端口使用,分割
		lp := make([]int64, 0)
		for _, port := range strings.Split(args[1], ",") {
			p, err := strconv.ParseInt(port, 10, 64)
			if err != nil {
				continue
			}
			lp = append(lp, p)
		}

		ports, err := funcs.ListenTcpPortMetrics(lp...)
		if err != nil {
			//todo 获得内存失败咋个处理
			log.Println(err)
			return
		}
		pkt := src.NewPkt()
		pkt.Id = g.PortListen
		pkt.Data, err = g.EncodeData(ports)
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

func (m PortService) Stop([]string) error {
	if portId > 0 {
		src.CancelTimer(portId)
	}
	return nil
}

func (m PortService) Restart(args []string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m PortService) Status([]string) bool {
	if portId > 0 {
		return true
	}
	return false
}
