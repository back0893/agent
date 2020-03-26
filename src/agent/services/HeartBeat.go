package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"github.com/back0893/goTcp/utils"
	"log"
	"time"
)

/**
心跳不能被中控服务器说控制
因为这个是一个必要的服务
*/

type HeartBeatService struct {
	CurrentStatus int   //当前配置的服务状态
	timerId       int64 //定时上报的定时器id
}

func (m *HeartBeatService) GetCurrentStatus() int {
	return m.CurrentStatus
}

func (m *HeartBeatService) SetCurrentStatus(status int) {
	m.CurrentStatus = status
}
func (m *HeartBeatService) Watcher() {
	//因为心跳是一个必要的服务,所以一旦心跳停止
	//需要重启,心跳服务
	run := m.Status(nil)
	if run == false {
		m.Start(map[string]string{})
	}
}

func (m *HeartBeatService) Upload(args map[string]string) {
	//这里讲发送分开
	//因为watcher方法,应该只能关注当前服务的状态
	//至于服务是如何上报,何时上报应该是由服务自己决定
	if m.timerId != 0 {
		src.CancelTimer(m.timerId)
	}
	interval := g.GetInterval(args, 5)
	m.timerId = src.AddTimer(interval*time.Second, func() {
		pkt := src.NewPkt()
		pkt.Id = g.PING
		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		if err := a.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	})
}
func NewHeartBeatService() *HeartBeatService {
	s := &HeartBeatService{
		CurrentStatus: 1,
	}
	s.Upload(map[string]string{})
	return s
}

func (m *HeartBeatService) Action(action string, args map[string]string) {
	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	switch action {
	case "start":
		m.Start(args)
	case "stop":
		m.Stop(args)
	case "restart":
		m.Restart(args)
	case "status":
		m.Status(args)
	case "interval":
		m.Upload(args)
	}
	pkt.Data = []byte("!启动心跳!")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m *HeartBeatService) Start(args map[string]string) error {
	if m.Status(nil) {
		return nil
	}
	m.CurrentStatus = 1
	return nil
}

func (m *HeartBeatService) Stop(map[string]string) error {
	m.CurrentStatus = 0
	return nil
}

func (m *HeartBeatService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HeartBeatService) Status(map[string]string) bool {
	return m.CurrentStatus == 1
}

func (m *HeartBeatService) Cancel() {
	src.CancelTimer(m.timerId)
}
