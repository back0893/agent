package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"github.com/back0893/goTcp/utils"
	"log"
)

type PortService struct {
	Ports []int64
}

func NewPortService(port []int64) *PortService {
	s := &PortService{
		Ports: port,
	}
	return s
}
func (m *PortService) SetPorts(ports []int64) {
	m.Ports = ports
}
func (m PortService) Info() {
	info := model.NewServiceResponse(g.PortListen, 1)
	portListen := m.Ports
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
