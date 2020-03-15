package model

type Service struct {
	Service string
	Action  string
	Args    map[string]string
}

func NewService(service, action string, args map[string]string) *Service {
	return &Service{
		Service: service,
		Action:  action,
		Args:    args,
	}
}
