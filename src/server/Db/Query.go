package Db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

/**
只需要实现4个方法的sql操作,就行,简单.
*/
const (
	INSERTSQL = "insert into %table% %field% values %value%"
)

type Query struct {
	db    *sql.DB
	table string
}

func sKv(value reflect.Value) (keys, values []string) {
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		vf := value.Field(i)
		//不可导出
		if tf.Anonymous {
			continue
		}
		//非法零值
		if !vf.IsValid() || reflect.DeepEqual(vf.Interface(), reflect.Zero(vf.Type()).Interface()) {
			continue
		}
		//获得指针值
		for vf.Type().Kind() == reflect.Ptr {
			vf = vf.Elem()
		}
		//如果是组合后的struct,那么需要重复调用获得嵌套的struct值
		if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
			cKeys, cValues := sKv(vf)
			keys = append(keys, cKeys...)
			values = append(values, cValues...)
		} else {
			//依据json 获得tag的ksy忽略无tag的字段
			key := strings.Split(tf.Tag.Get("json"), ",")[0]
			if key == "" {
				continue
			}
			value := format(vf)
			if value != "" {
				keys = append(keys, fmt.Sprintf("`%s`", key))
				values = append(values, value)
			}
		}
	}
	return
}
func format(value reflect.Value) string {
	//格式化输出
	if t, ok := value.Interface().(time.Time); ok {
		return t.Format("2006-01-02 15:04:05")
	}
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf(`'%s'`, value.Interface())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf(`%d`, value.Interface())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`%f`, value.Interface())
	case reflect.Slice, reflect.Array:
		var values []string
		for i := 0; i < value.Len(); i++ {
			values = append(values, format(value.Index(i)))
		}
		return fmt.Sprintf(`(%s)`, strings.Join(values, ","))
	case reflect.Interface:
		return format(value.Elem())
	default:
		return ""
	}
}
func mKv(value reflect.Value) (keys, values []string, err error) {
	itor := value.MapRange()
	for itor.Next() {
		if itor.Key().Kind() != reflect.String {
			return nil, nil, errors.New("map的key只能是string")
		}
		value := format(itor.Value())
		if value != "" {
			keys = append(keys, itor.Key().Interface().(string))
			values = append(values, value)
		}
	}
	return
}
func (query *Query) Insert(data interface{}) (int64, error) {
	if query.table == "" {
		return 0, errors.New("table为空")
	}
	var keys, values []string
	var err error
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
			_, val := sKv(sv)
			values = append(values, val...)
		}
	case reflect.Map:
		keys, values, err = mKv(v)
		if err != nil {
			return 0, err
		}
	default:
		return 0, errors.New("新增接受值错误")
	}

	kl := len(keys)
	vl := len(values)
	if kl == 0 || vl == 0 {
		return 0, errors.New("没有数据输入")
	}

	var insertValue string
	//如果是slice,kl一定大于vl
	if kl < vl {
		var tmpValues []string
		for kl <= vl {
			if kl%(len(keys)) == 0 {
				tmpValues = append(tmpValues, fmt.Sprintf("(%s)", strings.Join(values[kl-len(keys):kl], ",")))
			}
			kl++
		}
		insertValue = strings.Join(tmpValues, ",")
	} else {
		insertValue = fmt.Sprintf("(%s)", strings.Join(values, ","))
	}
	field := fmt.Sprintf("(%s)", strings.Join(keys, ","))
	replacer := strings.NewReplacer("%table%", query.table, "%field%", field, "%value%", insertValue)
	realSql := replacer.Replace(INSERTSQL)
	fmt.Println(realSql)
	result, err := query.db.Exec(realSql)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
func Table(db *sql.DB, tableName string) *Query {
	return &Query{
		db:    db,
		table: tableName,
	}
}
