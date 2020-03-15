package model

type ServiceStatus struct {
	Status int         // 当前的状态
	Data   interface{} //当前服务的额外值
}

type ServiceStatusList struct {
	Pid      int
	Services map[string]*ServiceStatus
}
