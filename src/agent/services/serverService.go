package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
)

/**
依据数据库,内存,cpu是一起
*/
type ServerService struct {
}

func NewServerService() *ServerService {
	return &ServerService{}
}
func (s *ServerService) Info() {
	info := model.NewServiceResponse(g.BaseServerInfo, 1)
	//获得cpu信息
	if funcs.CpuPrepared() == false {
		return
	}
	cpuStatus := funcs.CpuMetrics()
	//cpu的个数
	cpuNum := funcs.CpuNum()
	//cpu的基础频率
	cpuMhz, err := funcs.CpuMHz()
	if err != nil {
		cpuMhz = "0"
	}
	//获得内存
	mem, err := funcs.MemMetrics()
	if err != nil {
		fmt.Print(err)
		return
	}

	//获得负载
	loadAvg, err := funcs.LoadAvgMetrics()
	if err != nil {
		fmt.Print(err)
		return
	}

	data, err := g.EncodeData(cpuStatus, mem, loadAvg, cpuNum, cpuMhz)
	if err != nil {
		fmt.Print(err)
		return
	}
	info.Info = data
	pkt := src.ServiceResponsePkt(info)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err = a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
