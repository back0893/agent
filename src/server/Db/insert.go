package Db

import (
	"fmt"
	"reflect"
	"strings"
)

/**
只需要实现4个方法的sql操作,就行,简单.
*/
type A struct {
}
type SqlBuilder struct {
	table     string
	where     []string
	whereArgs []interface{}
}

func (builder *SqlBuilder) Table(table string) {
	builder.table = table
}
func (builder *SqlBuilder) Where(where string, args ...interface{}) {
	builder.where = append(builder.where, where)
	builder.whereArgs = append(builder.whereArgs, args...)
}
func (builder *SqlBuilder) Insert(t interface{}) {
	rv := reflect.ValueOf(t)
	ty := rv.Type()
	insert := []string{}
	for i := 0; i < ty.NumField(); i++ {
		insert = append(insert, ty.Field(i).Name)
	}

	sql := fmt.Sprintf("insert %s %s values %data%", builder.table, strings.Join(insert, ","))

}
func (builder *SqlBuilder) Query() {

}
func (builder *SqlBuilder) Update(t interface{}) {

}
func (builder *SqlBuilder) Delete() {

}
