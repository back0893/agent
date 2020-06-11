package model

type Auth struct {
	Username string
	Password string
	Id       int
}

//通常回应结构体

type Response struct {
	Id int32
}
