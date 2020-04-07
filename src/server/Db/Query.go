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
	UPDATESQL = "update %table% set %update% %where%"
	DELSQL    = "delete from %table%  %where%"
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
	errs   []string
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
func (q *Query) Insert(data interface{}) (int64, error) {
	if q.table == "" {
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
	replacer := strings.NewReplacer("%table%", q.table, "%field%", field, "%value%", insertValue)
	realSql := replacer.Replace(INSERTSQL)
	fmt.Println(realSql)
	result, err := q.db.Exec(realSql, values...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (q *Query) Where(w interface{}, args ...interface{}) *Query {
	if where, err := where("and", w, args...); err != nil {
		q.errs = append(q.errs, err.Error())
	} else {
		q.wheres = append(q.wheres, where)
	}
	return q
}
func (q *Query) WhereIn(w string, args ...interface{}) *Query {
	argsReapt := strings.Trim(strings.Repeat(",?", len(args)), ",")
	w = fmt.Sprintf("%s in (%s)", w, argsReapt)
	if where, err := where("and", w, args...); err != nil {
		q.errs = append(q.errs, err.Error())
	} else {
		q.wheres = append(q.wheres, where)
	}
	return q
}
func (q *Query) WhereOr(w interface{}, args ...interface{}) *Query {
	if where, err := where("or", w, args); err != nil {
		q.errs = append(q.errs, err.Error())
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
func (q *Query) where() (string, []interface{}) {
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
	return where, args
}
func (q *Query) Select(dest interface{}) error {
	if len(q.errs) > 0 {
		return errors.New(strings.Join(q.errs, "\n"))
	}
	v := reflect.ValueOf(dest)
	t := v.Type()
	if t.Kind() != reflect.Ptr {
		return errors.New("只能传递进入指针")
	}
	if !v.Elem().CanAddr() {
		return errors.New("不能正确的获得地址")
	}

	t = t.Elem()
	v = v.Elem()

	//如果没有field没有值,从struct取值
	if len(q.field) == 0 {
		switch t.Kind() {
		case reflect.Struct:
			if t.Name() != "Time" {
				q.field = sk(v)
			}
		case reflect.Slice:
			//如果是切面就需要取出其中的一个
			t := t.Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Struct {
				if t.Name() != "Time" {
					q.field = sk(reflect.Zero(t))
				}
			}
		}
	}
	if len(q.field) == 0 {
		return errors.New("查询的字段不能为空")
	}
	if t.Kind() != reflect.Slice {
		q.Limit(1)
	}
	//todo
	where, args := q.where()
	replacer := strings.NewReplacer("%field%", strings.Join(q.field, ","), "%table%", q.table, "%where%", where, "%order%", q.order, "%limit%", q.limit, "%offset%", q.offset)
	tmpSql := replacer.Replace(SELECTSQL)
	rows, err := q.db.Query(tmpSql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	switch t.Kind() {
	case reflect.Slice:
		dt := t.Elem()
		for dt.Kind() == reflect.Ptr {
			dt = dt.Elem()
		}
		sl := reflect.MakeSlice(t, 0, 0)
		for rows.Next() {
			var destination reflect.Value
			if dt.Kind() == reflect.Map {
				destination, err = q.setMap(rows, dt)
			} else {
				destination, err = q.setElem(rows, dt)
			}
			if err != nil {
				return err
			}
			switch t.Elem().Kind() {
			case reflect.Ptr, reflect.Map:
				sl = reflect.Append(sl, destination)
			default:
				sl = reflect.Append(sl, destination.Elem())
			}
		}
		v.Set(sl)
		return nil
	case reflect.Map:
		for rows.Next() {
			m, err := q.setMap(rows, t)
			if err != nil {
				return err
			}
			v.Set(m)
		}
		return nil
	default:
		for rows.Next() {
			destination, err := q.setElem(rows, t)
			if err != nil {
				return err
			}
			v.Set(destination.Elem())
		}
	}
	return nil
}
func (q *Query) Field(columns ...string) *Query {
	q.field = append(q.field, columns...)
	return q
}
func address(dest reflect.Value, columns []string) []interface{} {
	dest = dest.Elem()
	t := dest.Type()
	addrs := make([]interface{}, 0)
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			vf := dest.Field(i)
			tf := t.Field(i)
			if tf.Anonymous {
				continue
			}
			for vf.Kind() == reflect.Ptr {
				vf = vf.Elem()
			}
			if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
				nvf := reflect.New(vf.Type())
				vf.Set(nvf.Elem())
				addrs = append(addrs, address(nvf, columns)...)
				continue
			}
			column := strings.Split(tf.Tag.Get("json"), ",")[0]
			if column == "" {
				continue
			}
			for _, col := range columns {
				if col == column {
					addrs = append(addrs, vf.Addr().Interface())
					break
				}
			}
		}
	default:
		addrs = append(addrs, dest.Addr().Interface())
	}
	return addrs
}

//因为map不能用new生成,所以只能用一个方法来生成
func (q *Query) setMap(rows *sql.Rows, t reflect.Type) (reflect.Value, error) {
	if t.Elem().Kind() != reflect.Interface {
		return reflect.ValueOf(nil), errors.New("map的值只能是interface")
	}
	m := reflect.MakeMap(t)
	addrs := make([]interface{}, len(q.field))
	for idx := range q.field {
		addrs[idx] = new(interface{})
	}
	if err := rows.Scan(addrs...); err != nil {
		return reflect.ValueOf(nil), err
	}
	for idx, column := range q.field {
		m.SetMapIndex(reflect.ValueOf(column), reflect.ValueOf(addrs[idx]).Elem().Elem())
	}
	return m, nil
}
func (q *Query) setElem(rows *sql.Rows, t reflect.Type) (reflect.Value, error) {
	addrsErr := errors.New("不能匹配")
	dest := reflect.New(t)
	addrs := address(dest, q.field)
	if len(q.field) != len(addrs) {
		return reflect.ValueOf(nil), addrsErr
	}
	if err := rows.Scan(addrs...); err != nil {
		return reflect.ValueOf(nil), err
	}
	return dest, nil
}

func (q *Query) Update(src interface{}) (int64, error) {
	if len(q.errs) != 0 {
		return 0, errors.New(strings.Join(q.errs, "\n"))
	}
	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	var updated string
	var keys []string
	var values []interface{}
	switch v.Kind() {
	case reflect.String:
		updated = v.Interface().(string)
	case reflect.Struct:
		keys, values = sKv(v)
	case reflect.Map:
		keys, values, _ = mKv(v)
	default:
		return 0, errors.New("不支持的类型")
	}
	if updated == "" {
		if len(keys) != len(values) {
			return 0, errors.New("更新的字段无法对应")
		}
		kvs := make([]string, 0)
		for _, key := range keys {
			kvs = append(kvs, fmt.Sprintf("%s=?", key))
		}
		updated = strings.Join(kvs, ",")
	}
	where, args := q.where()
	values = append(values, args...)
	replacer := strings.NewReplacer("%table%", q.table, "%update%", updated, "%where%", where)
	tmpSql := replacer.Replace(UPDATESQL)
	result, err := q.db.Exec(tmpSql, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func (q *Query) Delete() (int64, error) {
	if len(q.errs) != 0 {
		return 0, errors.New(strings.Join(q.errs, "\n"))
	}
	if len(q.wheres) == 0 {
		return 0, errors.New("删除条件不能为空")
	}
	where, args := q.where()
	replacer := strings.NewReplacer("%table%", q.table, "%where%", where)
	tmpSql := replacer.Replace(DELSQL)
	result, err := q.db.Exec(tmpSql, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
func sk(value reflect.Value) []string {
	var keys []string
	t := value.Type()
	for i := 0; i < t.NumField(); i++ {
		tf := t.Field(i)
		vf := value.Field(i)
		if tf.Anonymous {
			continue
		}
		for vf.Kind() == reflect.Ptr {
			vf = vf.Elem()
		}

		if vf.Kind() == reflect.Struct && tf.Type.Name() != "Time" {
			keys = append(keys, sk(vf)...)
			continue
		}
		key := strings.Split(tf.Tag.Get("json"), ",")[0]
		if key == "" {
			continue
		}
		keys = append(keys, key)
	}
	return keys
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
