package Db

import (
	"database/sql"
	"errors"
	"golang.org/x/text/cases"
	"reflect"
)

/**
只需要实现4个方法的sql操作,就行,简单.
*/

type Query struct {
	db    *sql.DB
	table string
}

func sKv(value reflect.Value) (keys, values []string) {

}
func mKv(value reflect.Value) (keys, values []string) {

}
func (query *Query) Insert(data interface{}) (int, error) {
	var keys, values []string
	v := reflect.ValueOf(data)
	//如果data是一个指针,获得指针的值
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		keys, values = sKv(v)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			sv := v.Index(i)
			for sv.Kind() == reflect.Ptr || sv.Kind() == reflect.Interface {
				sv = sv.Elem()
			}
			if sv.Kind() != reflect.Struct {
				return 0, errors.New("slice 只接受struct")
			}
			if len(keys) == 0 {
				keys, values = sKv(sv)
				continue
			}
			_, val := sKv()
			values = append(values, val...)
		}
	case reflect.Map:
		keys, values = mKv(v)
	default:
		return 0, errors.New("新增接受值错误")
	}
	return 0, nil
}
func Table(db *sql.DB, tableName string) *Query {
	return &Query{
		db:    db,
		table: tableName,
	}
}
