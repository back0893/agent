package services

import (
	"agent/src/agent/funcs"
	g2 "agent/src/agent/g"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"github.com/back0893/goTcp/utils"
	"log"
)

type PortService struct {
}

func NewPortService() *PortService {
	return &PortService{}
}

func (m PortService) Info() {
	info := model.NewServiceResponse(g.PortListen, 1)
	portListen := g2.GetPortListen()
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
	pkt := g.ServiceResponsePkt(info)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err = a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
