package model

type Service struct {
	Service int
	Status  int
	Info    []byte //每个服务返回的值,响应的服务自己处理
}

func NewServiceResponse(service, status int) *Service {
	return &Service{
		Service: service,
		Status:  status,
	}
}
