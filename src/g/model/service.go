package model

type Service struct {
	Service int
	Action  string
	Args    map[string]string
}

func NewService(service int, action string, args map[string]string) *Service {
	return &Service{
		Service: service,
		Action:  action,
		Args:    args,
	}
}

type ServiceResponse struct {
	Service int
	Status  int
	Info    []byte //每个服务返回的值,响应的服务自己处理
}

func NewServiceResponse(service, status int) *ServiceResponse {
	return &ServiceResponse{
		Service: service,
		Status:  status,
	}
}
