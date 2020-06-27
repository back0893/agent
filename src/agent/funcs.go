package agent

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/agent/services"
	"github.com/back0893/goTcp/utils"
	"time"
)

type FuncsAndInterval struct {
	Fs       []iface.IService
	Interval int
}

var Mappers []*FuncsAndInterval

func BuildMappers() {
	interval := utils.GlobalConfig.GetInt("heartBeat")
	Mappers = []*FuncsAndInterval{
		{
			Fs: []iface.IService{
				services.NewHeartBeatService(),
			},
			Interval: interval,
		},
		{
			Fs: []iface.IService{
				services.NewPortService(),
			},
			Interval: 120,
		},
		{
			Fs: []iface.IService{
				services.NewServerService(),
				services.NewHHDService(),
			},
			Interval: 60,
		},
	}
}
func Collect() {
	for _, v := range Mappers {
		src.AddTimer(time.Duration(v.Interval)*time.Second, func(v *FuncsAndInterval) func() {
			//返回一个闭包,因为在执行时v已经被改变..
			return func() {
				for _, service := range v.Fs {
					service.Info()
				}
			}
		}(v))
	}
}
