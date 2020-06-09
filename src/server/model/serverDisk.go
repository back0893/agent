package model

type ServerDisk struct {
	Id       int64
	Name     string
	Gb       string
	ServerId int `db:"server_id"`
}
