package agent

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/agent/services"
	"agent/src/g"
	"agent/src/g/model"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type ServicesList struct {
	services  map[int]iface.IService
	taskQueue *src.TaskQueue
}

func NewServicesList() *ServicesList {
	//心跳的服务默认存在
	return &ServicesList{
		services:  map[int]iface.IService{},
		taskQueue: src.NewTaskQueue(),
	}

}
func (sl *ServicesList) AddService(name int, s iface.IService) {
	sl.services[name] = s
}
func (sl *ServicesList) GetService(name int) (service iface.IService, ok bool) {
	service, ok = sl.services[name]
	return
}
func (sl *ServicesList) GetServices() map[int]iface.IService {
	return sl.services
}
func (sl *ServicesList) CancelService(name int) {
	if service, ok := sl.services[name]; ok {
		service.Cancel()
	}
	delete(sl.services, name)
}
func (sl *ServicesList) CancelAll() {
	for name, _ := range sl.services {
		service := sl.services[name]
		service.Cancel()
		delete(sl.services, name)
	}
}
func (sl *ServicesList) BaseService() {
	//心跳必须存在
	sl.AddService(g.PING, services.NewHeartBeatService())
	sl.AddService(g.BaseServerInfo, services.NewServerService(1))
	sl.AddService(g.HHD, services.NewHHDService(1))
	//sl.AddService(g.PortListen,services.NewPortService(1))
}
func (sl *ServicesList) WakeUp() map[int]int {
	path := g.GetRuntimePath()
	t := make(map[int]int)
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, "services"))
	if err != nil {
		return t
	}
	if err := g.DecodeData(data, &t); err != nil {
		return t
	}
	return t
}
func (sl *ServicesList) Sleep() error {
	t := make(map[int]int)
	for name, s := range sl.GetServices() {
		t[name] = s.GetCurrentStatus()
	}
	path := g.GetRuntimePath()
	data, err := g.EncodeData(t)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s/%s", path, "services"), data, 0644); err != nil {
		return err
	}
	return nil
}

func (sl *ServicesList) NewService(name, status int) (iface.IService, error) {
	var service iface.IService
	switch name {
	case g.REDISSERVICE:
		service = services.NewRedisService(status)
	default:
		return nil, errors.New("服务还未被实现")
	}
	return service, nil
}

/**
和从中控服务器下发的启动service同步
*/
func (sl *ServicesList) Sync(data []byte) {
	ss := make(map[int]int, 0)
	if err := g.DecodeData(data, &ss); err != nil {
		fmt.Println(err)
		return
	}
	store := sl.WakeUp()
	for name, status := range ss {
		t, ok := store[name]
		if ok {
			status = t
		}
		if service, err := sl.NewService(name, status); err == nil {
			sl.AddService(name, service)
		}
	}
}

/**
执行对应的服务动作
*/
func (sl *ServicesList) RunServiceAction() {
	//读取taskQueue,执行相应的操作
	for {
		task := sl.taskQueue.Pop()
		fmt.Print(task.Service)

		//如果是一个取消服务的动作
		if strings.ToLower(task.Action) == "cancel" {
			sl.CancelService(task.Service)
			return
		}
		service, ok := sl.GetService(task.Service)
		if ok == false {
			//todo 回应服务未启动?
			fmt.Println("当前服务临时启动...")
			var err error
			service, err = sl.NewService(task.Service, 1)
			if err != nil {
				fmt.Println("启动服务失败...")
				return
			}
		}
		fmt.Printf("%p\n", service)
		service.Action(task.Action, task.Args)
	}
}

/**
新增一个服务的动作
*/

func (sl *ServicesList) AddServiceAction(task *model.Service) {
	sl.taskQueue.Push(task)
}

/**
循环监控所有服务的状态
*/
func (sl *ServicesList) Listen() {
	for _, service := range sl.GetServices() {
		service.Watcher()
	}
}
