package model

type Service struct {
	Service string
	Action  string
	Args    []string
}

func NewService(service, action string, args ...string) *Service {
	return &Service{
		Service: service,
		Action:  action,
		Args:    args,
	}
}
