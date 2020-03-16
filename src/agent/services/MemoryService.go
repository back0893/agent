package services

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/iface"
	"agent/src/g"
	"errors"
	"github.com/back0893/goTcp/utils"
	"log"
	"strconv"
	"time"
)

//一个私有的全局变量
var memId int64

type MemoryService struct {
}

func NewMemoryService() *MemoryService {
	return &MemoryService{}
}

func (m *MemoryService) Action(action string, args map[string]string) {
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
	pkt := src.NewPkt()
	pkt.Id = g.ServiceResponse
	pkt.Data = []byte("启动memory")

	a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	err := a.GetCon().Write(pkt)

	if err != nil {
		//todo 发送失败..应该有后续操作
	}
}

func (m MemoryService) Start(args map[string]string) error {
	if m.Status(args) {
		return errors.New("service已经启动")
	}
	var num = 60
	if len(args) > 0 {
		n, err := strconv.Atoi(args["interval"])
		if err == nil {
			num = n
		}
	}
	memId = src.AddTimer(time.Duration(num)*time.Second, func() {
		memory, err := funcs.MemMetrics()
		if err != nil {
			//todo 获得内存失败咋个处理
			log.Println(err)
			return
		}
		pkt := src.NewPkt()
		pkt.Id = g.MEM
		pkt.Data, err = g.EncodeData(memory)
		if err != nil {
			log.Println(err)
			return
		}
		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		err = a.GetCon().Write(pkt)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return nil
}

func (m MemoryService) Stop(map[string]string) error {
	if memId > 0 {
		src.CancelTimer(memId)
	}
	return nil
}

func (m MemoryService) Restart(args map[string]string) error {
	if err := m.Stop(args); err != nil {
		return err
	}
	if err := m.Start(args); err != nil {
		return err
	}
	return nil
}

func (m MemoryService) Status(map[string]string) bool {
	if memId > 0 {
		return true
	}
	return false
}
