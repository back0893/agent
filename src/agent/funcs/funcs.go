package funcs

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
var portChan chan int64

func BuildMappers() {
	interval := utils.GlobalConfig.GetInt("device.interval")
	portService := services.NewPortService([]int64{})
	portChan = make(chan int64, 10)
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
	go func() {
		for port := range portChan {
			portService.Ports = append(portService.Ports, port)
		}
	}()
}
func AppendPorts(ports ...int64) {
	for _, port := range ports {
		portChan <- port
	}
}
