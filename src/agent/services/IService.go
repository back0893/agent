package services

/**
控制外部资源
对于外部资源,必须有启动,停止和重启,状态的查询
*/
type IService interface {
	Start()
	Stop()
	Restart()
	Status()
}
