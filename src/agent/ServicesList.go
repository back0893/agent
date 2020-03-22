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

type service struct {
	Name          string
	CurrentStatus string
}
type ServicesList struct {
	services  map[string]iface.IService
	taskQueue *src.TaskQueue
}

func NewServicesList() *ServicesList {
	//心跳的服务默认存在
	return &ServicesList{
		services:  map[string]iface.IService{},
		taskQueue: src.NewTaskQueue(),
	}

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
func (sl *ServicesList) CancelService(name string) {
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

func (sl *ServicesList) WakeUp() error {
	path := g.GetRuntimePath()
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, "services"))
	if err != nil {
		return err
	}
	t := make([]service, 0)
	if err := g.DecodeData(data, &t); err != nil {
		return err
	}
	for _, se := range t {
		service, err := sl.NewService(se.Name)
		if err != nil {
			continue
		}
		service.SetCurrentStatus(se.CurrentStatus)
		sl.AddService(se.Name, service)
		fmt.Println(se.Name, se.CurrentStatus)
	}
	//唤醒后心跳必须存在
	sl.CancelService("heart")
	sl.AddService("heart", services.NewHeartBeatService())
	return nil
}
func (sl *ServicesList) Sleep() error {
	t := make([]service, 0)
	for name, s := range sl.GetServices() {
		t = append(t, service{
			Name:          name,
			CurrentStatus: s.GetCurrentStatus(),
		})
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
	if err := g.DecodeData(data, &ss); err != nil {
		fmt.Println(err)
		return
	}
	//中控下发的服务+本地已经存在的服务
	for _, name := range ss {
		if _, ok := sl.services[name]; ok == false {
			fmt.Println("-------------------new service")
			if service, err := sl.NewService(name); err == nil {
				sl.AddService(name, service)
			}
		}
	}
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
		fmt.Print(task.Service)

		//如果是一个取消服务的动作
		if strings.ToLower(task.Action) == "cancel" {
			sl.CancelService(task.Service)
			return
		}

		service, ok = sl.GetService(task.Service)
		fmt.Println("get=====>", ok)
		fmt.Printf("%p\n", service)
		fmt.Printf("%p\n", sl.services["redis"])
		if ok == false {
			service, err = sl.NewService(task.Service)
			if err != nil {
				//todo 服务未被实现
				continue
			}
			sl.AddService(task.Service, service)
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
