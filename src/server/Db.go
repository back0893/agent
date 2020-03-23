package server

import (
	"github.com/back0893/goTcp/utils"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

func init() {
	dsn := mysql.Config{
		User:   utils.GlobalConfig.GetString("db.username"),
		Passwd: utils.GlobalConfig.GetString("db.password"),
		Net:    utils.GlobalConfig.GetString("db.net"),
		Addr:   utils.GlobalConfig.GetString("db.addr"),
		DBName: utils.GlobalConfig.GetString("db.dbname"),
		Params: utils.GlobalConfig.GetStringMapString("db.params"),
	}

	var err error
	db, err = sqlx.Open("mysql", dsn.FormatDSN())
	if err != nil {
		panic(err)
	}

}
