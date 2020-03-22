package iface

/**
控制外部资源
对于外部资源,必须有启动,停止和重启,状态的查询
*/
type IService interface {
	Start(map[string]string) error
	Stop(map[string]string) error
	Restart(map[string]string) error
	Status(map[string]string) bool                //查询服务的运行状态
	Action(action string, args map[string]string) //用来执行相应的动作
	Watcher()                                     //监控运行状态
	GetCurrentStatus() string
	SetCurrentStatus(string)
	Cancel() //取消服务时执行的任务
}
