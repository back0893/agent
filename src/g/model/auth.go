package model

type Auth struct {
	Username string
	Password string
	Id       int
}

//AuthResponse 认真回应
type AuthResponse struct {
	//是否成功
	Status bool
}
