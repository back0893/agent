package services

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	"errors"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
	"strconv"
	"time"
)

//一个私有的全局变量
var heartId int64

type HeartBeatService struct {
}

func NewHeartBeatService() *HeartBeatService {
	return &HeartBeatService{}
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
	}
	pkt.Data = []byte("!启动心跳!")
	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)
	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m HeartBeatService) Start(args map[string]string) error {
	//如果已经启动,,,不能重复启动
	if m.Status(args) {
		return errors.New("service已经启动")
	}
	var num = 10
	if len(args) > 0 {
		n, err := strconv.Atoi(args["interval"])
		if err == nil {
			num = n
		}
	}
	fmt.Println(num)
	heartId = src.AddTimer(time.Duration(num)*time.Second, func() {
		fmt.Println("ping")
		pkt := src.NewPkt()
		pkt.Id = g.PING
		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		if err := a.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	})
	return nil
}

func (m HeartBeatService) Stop(map[string]string) error {
	if heartId > 0 {
		src.CancelTimer(heartId)
	}
	heartId = 0
	return nil
}

func (m HeartBeatService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m HeartBeatService) Status(map[string]string) bool {
	if heartId > 0 {
		return true
	}
	return false
}
