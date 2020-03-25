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
