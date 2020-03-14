package iface

/**
控制外部资源
对于外部资源,必须有启动,停止和重启,状态的查询
*/
type IService interface {
	Start(map[string]string) error
	Stop(map[string]string) error
	Restart(map[string]string) error
	Status(map[string]string) bool //查询运营状态
	Action(action string, args map[string]string)
}
