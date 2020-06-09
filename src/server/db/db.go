package Db

import (
	"agent/src/g"
	"github.com/back0893/goTcp/utils"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

type dbCons struct {
	dbs  map[string]*sqlx.DB
	lock sync.RWMutex
}

func (db *dbCons) Get(name string) (*sqlx.DB, bool) {
	db.lock.RLock()
	defer db.lock.RUnlock()
	con, ok := db.dbs[name]
	if ok == true {
		return con, true
	}
	if con, err := newConnection(name); err != nil {
		return nil, false
	} else {
		db.dbs[name] = con
		return con, true
	}
}

var (
	DbConnections *dbCons
)

func init() {
	DbConnections = &dbCons{
		dbs: make(map[string]*sqlx.DB),
	}
}

func newConnection(name string) (*sqlx.DB, error) {
	dsn := mysql.Config{
		User:                 utils.GlobalConfig.GetString("db.username"),
		Passwd:               utils.GlobalConfig.GetString("db.password"),
		Net:                  utils.GlobalConfig.GetString("db.net"),
		Addr:                 utils.GlobalConfig.GetString("db.addr"),
		DBName:               utils.GlobalConfig.GetString("db.dbname"),
		Params:               utils.GlobalConfig.GetStringMapString("db.params"),
		AllowNativePasswords: true,
	}

	dsn.Loc = g.CSTLocation()
	db, err := sqlx.Open("mysql", dsn.FormatDSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
