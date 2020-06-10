package agent

import (
	"agent/src/agent/iface"
	"agent/src/agent/services"
	"github.com/back0893/goTcp/utils"
)

type FuncsAndInterval struct {
	Fs       []iface.IService
	Interval int
}

var Mappers []FuncsAndInterval

func BuildMappers() {
	interval := utils.GlobalConfig.GetInt("heartBeat")
	portService := services.NewPortService([]int64{})
	Mappers = []FuncsAndInterval{
		{
			Fs: []iface.IService{
				services.NewServerService(),
				services.NewHeartBeatService(),
				services.NewHHDService(),
			},
			Interval: interval,
		},
		{
			Fs: []iface.IService{
				portService,
			},
			Interval: interval,
		},
	}

}
