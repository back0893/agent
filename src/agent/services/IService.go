package services

/**
控制外部资源
对于外部资源,必须有启动,停止和重启,状态的查询
*/
type IService interface {
	Start() error
	Stop() error
	Restart() error
	Status() bool //查询运营状态
}
