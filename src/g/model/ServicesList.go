package model

type ServicesList struct {
	services map[string]*Service
}

func NewServicesList() *ServicesList {
	return &ServicesList{services: map[string]*Service{}}
}
func (sl *ServicesList) AddService(name string, s *Service) {
	sl.services[name] = s
}
func (sl *ServicesList) GetService(name string) (service *Service, ok bool) {
	service, ok = sl.services[name]
	return
}
func (sl *ServicesList) GetServices() map[string]*Service {
	return sl.services
}
