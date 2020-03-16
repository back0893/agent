package agent

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/agent/services"
	"agent/src/g"
	"errors"
	"io/ioutil"
)

type ServicesList struct {
	services  map[string]iface.IService
	agent     *Agent
	taskQueue *src.TaskQueue
}

func NewServicesList() *ServicesList {
	return &ServicesList{services: map[string]iface.IService{}}
}
func (sl *ServicesList) AddService(name string, s iface.IService) {
	sl.services[name] = s
}
func (sl *ServicesList) GetService(name string) (service iface.IService, ok bool) {
	service, ok = sl.services[name]
	return
}
func (sl *ServicesList) GetServices() map[string]iface.IService {
	return sl.services
}

func (sl *ServicesList) WakeUp() error {
	path := g.GetRuntimePath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if err := g.DecodeData(data, sl.services); err != nil {
		return err
	}
	return nil
}
func (sl *ServicesList) Sleep() error {
	path := g.GetRuntimePath()
	data, err := g.EncodeData(sl.services)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (sl *ServicesList) NewService(name string) (iface.IService, error) {
	var service iface.IService
	switch name {
	case "redis":
		service = services.NewRedisService()
	case "heart":
		service = services.NewHeartBeatService()
	case "loadavg":
		service = services.NewLoadAvgServiceService()
	case "memory":
		service = services.NewMemoryService()
	case "hhd":
		service = services.NewHHDService()
	case "port":
		service = services.NewPortService()
	case "cpu":
		service = services.NewCPUService()
	default:
		return nil, errors.New("服务还未被实现")
	}
	return service, nil
}

/**
和从中控服务器下发的启动service同步
*/
func (sl *ServicesList) Sync(data []byte) {
	ss := make([]string, 0)
	if err := g.DecodeData(data, ss); err != nil {
		return
	}
	sync := make(map[string]iface.IService)
	for _, name := range ss {
		if tmp, ok := sl.services[name]; ok == false {
			if service, err := sl.NewService(name); err == nil {
				sl.AddService(name, service)
			}
		} else {
			sync[name] = tmp
		}
	}
	sl.services = sync
}

/**
执行对应的服务动作
*/
func (sl *ServicesList) RunServiceAction() {
	//读取taskQueue,执行相应的操作
	for {
		var service iface.IService
		var ok bool
		var err error

		task := sl.taskQueue.Pop()
		service, ok = sl.GetService(task.Service)
		if ok == false {
			service, err = sl.NewService(task.Service)
			if err != nil {
				//todo 服务未被实现
				continue
			}
			sl.AddService(task.Service, service)
		}
		service.Action(task.Action, task.Args)
	}
}
