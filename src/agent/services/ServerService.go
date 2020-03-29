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
	"time"
)

/**
依据数据库,内存,cpu是一起
*/
type ServerService struct {
	timerId       int64
	CurrentStatus int
}

func NewServerService(status int) *ServerService {
	s := &ServerService{CurrentStatus: status}
	s.Upload(map[string]string{})
	return s
}

func (s *ServerService) Start(args map[string]string) error {
	if s.Status(args) {
		return nil
	}
	s.CurrentStatus = 1
	return nil
}

func (s *ServerService) Stop(map[string]string) error {
	s.CurrentStatus = 0
	return nil
}

func (s *ServerService) Restart(args map[string]string) error {
	if err := s.Stop(args); err != nil {
		return err
	}
	if err := s.Start(args); err != nil {
		return err
	}
	return nil
}

func (s *ServerService) Status(map[string]string) bool {
	return s.CurrentStatus == 1
}

func (s *ServerService) Action(action string, args map[string]string) {
	switch action {
	case "start":
		s.Start(args)
	case "stop":
		s.Stop(args)
	case "restart":
		s.Restart(args)
	case "status":
		s.info()
	}
	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	pkt.Data = []byte("操作完成")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (s *ServerService) Watcher() {
	run := s.Status(nil)
	if run == true && s.CurrentStatus == 0 {
		s.CurrentStatus = 1
	} else if s.CurrentStatus == 1 && run == false {
		s.Start(map[string]string{})
	}
}

func (s *ServerService) GetCurrentStatus() int {
	return s.CurrentStatus
}

func (s *ServerService) SetCurrentStatus(status int) {
	s.CurrentStatus = status
}

func (s *ServerService) Cancel() {
	src.CancelTimer(s.timerId)
}
func (s *ServerService) info() {
	//使用当前服务的状态作为默认值
	//如果启动失败,那么把状态修改成失败
	info := model.NewServiceResponse(g.BaseServerInfo, s.CurrentStatus)

	//获得cpu信息
	if funcs.CpuPrepared() == false {
		funcs.UpdateCpuStat()
	}
	err := funcs.UpdateCpuStat()
	if err != nil {
		info.Info = []byte("获得cpu失败")
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

func (s *ServerService) Upload(args map[string]string) {
	if s.timerId != 0 {
		src.CancelTimer(s.timerId)
	}
	interval := g.GetInterval(args, 60)
	s.timerId = src.AddTimer(interval*time.Second, s.info)
}
