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

type HHDService struct{}

func NewHHDService() *HHDService {
	return &HHDService{}
}
func (m HHDService) Info() {
	info := model.NewServiceResponse(g.HHD, 1)
	disks, err := funcs.DiskUseMetrics()
	if err != nil {
		log.Println(err)
		return
	}
	info.Info, err = g.EncodeData(disks)
	pkt := src.ServiceResponsePkt(info)
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err = a.GetCon().Write(pkt)
	if err != nil {
		log.Println(err)
		return
	}
}
