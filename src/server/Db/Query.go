package Db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**
只需要实现4个方法的sql操作,就行,简单.
*/
const (
	INSERTSQL = "insert into %table% %field% values %value%"
	SELECTSQL = "select %field% from %table% %where% %order% %limit% %offset%"
)

type Where struct {
	logic string        //连接方式
	where string        //where条件
	args  []interface{} //参数
}

func (w Where) GetWhere() string {
	return w.where
}
func (w Where) GetArgs() []interface{} {
	return w.args
}

type Query struct {
	db     *sql.DB
	table  string
	wheres []*Where
	field  []string
	limit  string
	offset string
	order  string
	errs   []error
}

func sKv(value reflect.Value) (keys []string, values []interface{}) {
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
			keys = append(keys, fmt.Sprintf("`%s`", key))
			values = append(values, vf.Interface())
		}
	}
	return
}
func mKv(value reflect.Value) (keys []string, values []interface{}, err error) {
	itor := value.MapRange()
	for itor.Next() {
		if itor.Key().Kind() != reflect.String {
			return nil, nil, errors.New("map的key只能是string")
		}
		keys = append(keys, itor.Key().Interface().(string))
		values = append(values, itor.Value())
	}
	return
}
func (query *Query) Insert(data interface{}) (int64, error) {
	if query.table == "" {
		return 0, errors.New("table为空")
	}
	var keys []string
	var values []interface{}
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
	if vl%kl == 0 {
		var tmpValues []string
		for i := vl / kl; i > 0; i-- {
			if kl%(len(keys)) == 0 {
				insertValue = fmt.Sprintf("(%s)", strings.Trim(strings.Repeat(",?", kl), ","))
				tmpValues = append(tmpValues, insertValue)
			}
		}
		insertValue = strings.Join(tmpValues, ",")
	} else {
		return 0, errors.New("插入长度不一致")
	}

	field := fmt.Sprintf("(%s)", strings.Join(keys, ","))
	replacer := strings.NewReplacer("%table%", query.table, "%field%", field, "%value%", insertValue)
	realSql := replacer.Replace(INSERTSQL)
	fmt.Println(realSql)
	result, err := query.db.Exec(realSql, values...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (q *Query) Where(w interface{}, args ...interface{}) *Query {
	if where, err := where("and", w, args...); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.wheres = append(q.wheres, where)
	}
	return q
}
func (q *Query) WhereOr(w interface{}, args ...interface{}) *Query {
	if where, err := where("or", w, args); err != nil {
		q.errs = append(q.errs, err)
	} else {
		q.wheres = append(q.wheres, where)
	}
	return q
}
func (q *Query) Limit(limit uint) *Query {
	q.limit = fmt.Sprintf("limit %d", limit)
	return q
}
func (q *Query) Offset(offset uint) *Query {
	q.offset = fmt.Sprintf("offset %d", offset)
	return q
}
func (q *Query) Order(ord string, asc bool) *Query {
	t := "asc"
	if asc == false {
		t = "desc"
	}
	q.order = fmt.Sprintf("order by %s %s", ord, t)
	return q
}

func (q *Query) Select() *sql.Row {
	where := ""
	args := make([]interface{}, 0)
	if len(q.wheres) > 0 {
		t := make([]string, len(q.wheres))
		for i, w := range q.wheres {
			t[i] = fmt.Sprintf("(%s)", w.GetWhere())
			args = append(args, w.GetArgs()...)
		}
		where = "where " + strings.Join(t, " and ")
	}
	replacer := strings.NewReplacer("%field%", "`id`,`name`", "%table%", q.table, "%where%", where, "%order%", q.order, "%limit%", q.limit, "%offset%", q.offset)
	sqlStr := replacer.Replace(SELECTSQL)
	return q.db.QueryRow(sqlStr, args...)
}
func newWhere(logic string, w string, args ...interface{}) *Where {
	return &Where{
		logic: logic,
		where: w,
		args:  args,
	}
}
func where(logic string, w interface{}, args ...interface{}) (*Where, error) {
	v := reflect.ValueOf(w)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return newWhere(logic, w.(string), args...), nil
	case reflect.Struct:
		keys, values := sKv(v)
		for i, _ := range keys {
			keys[i] = fmt.Sprintf("%s=?", keys[i])
		}
		return newWhere(logic, strings.Join(keys, ","), values...), nil
	case reflect.Map:
		keys, values, err := mKv(v)
		if err != nil {
			return nil, err
		}
		for i, _ := range keys {
			keys[i] = fmt.Sprintf("%s=?", keys[i])
		}
		return newWhere(logic, strings.Join(keys, ","), values...), nil
	default:
		return nil, errors.New("不支持的条件参数")
	}

}
func Table(db *sql.DB, tableName string) *Query {
	return &Query{
		db:    db,
		table: tableName,
	}
}
