package model

type Service struct {
	TemplateId int `db:"template_id"`
	Status     int
}
