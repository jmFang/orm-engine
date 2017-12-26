package entities

import (
	"bytes"
	"container/list"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//时间格式
const (
	timeFormate = "2006-01-02 15:04:05"
)

var sqlDB *sql.DB

func Open() *sql.DB {
	if sqlDB == nil {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
		if err != nil {
			fmt.Println(err.Error())
		}
		sqlDB = db
	}

	return sqlDB
}

type IBaseDao interface {
	Init()
	Save(data interface{}) error
	Update(data interface{}) error
	Find() (*list.List, error)
}

type BaseDao struct {
	EntityType    reflect.Type
	sqlDB         *sql.DB
	tableName     string            //表名
	pk            string            //主键
	columnToField map[string]string //字段名:属性名
	fieldToColumn map[string]string //属性名:字段名
}

//初始化
func (this *BaseDao) Init() {
	this.columnToField = make(map[string]string)
	this.fieldToColumn = make(map[string]string)

	types := this.EntityType

	for i := 0; i < types.NumField(); i++ {
		typ := types.Field(i)
		tag := typ.Tag

		if len(tag) > 0 {
			column := tag.Get("column")
			name := typ.Name
			this.columnToField[column] = name
			this.fieldToColumn[name] = column

			if len(tag.Get("table")) > 0 {
				this.tableName = tag.Get("table")
				this.pk = column
			}
		}
	}
}

//预处理插入sql
func (this *BaseDao) insertPrepareSQL() (fieldNames list.List, sql string) {
	names := new(bytes.Buffer)
	values := new(bytes.Buffer)

	i := 0

	for column, fieldName := range this.columnToField {

		if i != 0 {
			names.WriteString(",")
			values.WriteString(",")
		}
		fieldNames.PushBack(fieldName)
		names.WriteString(column)
		values.WriteString("?")
		i++
	}
	sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", this.tableName, names.String(), values.String())
	return
}

//增加单列
func (this *BaseDao) Save(data interface{}) error {
	columns, sql := this.insertPrepareSQL()

	stmt, err := Open().Prepare(sql)
	args := this.prepareValues(data, columns)
	fmt.Println(sql, " ", args)
	_, err = stmt.Exec(args...)
	if err != nil {
		panic(err.Error())
	}
	return err
}

//更新一个实体
func (this *BaseDao) Update(data interface{}) error {
	columns, sql := this.updatePrepareSQL()

	stmt, err := Open().Prepare(sql)
	args := this.prepareValues(data, columns)

	fmt.Println(sql, " ", args)
	_, err = stmt.Exec(args...)
	if err != nil {
		panic(err.Error())
	}
	return err
}

//实体转update sql语句
func (this *BaseDao) updatePrepareSQL() (fieldNames list.List, sql string) {
	//UPDATE 表名称 SET 列名称 = 新值 WHERE 列名称 = 某值
	sets := new(bytes.Buffer)

	i := 0

	for column, fieldName := range this.columnToField {
		if strings.EqualFold(column, this.pk) {
			continue
		}
		if i != 0 {
			sets.WriteString(",")
		}

		fieldNames.PushBack(fieldName)
		sets.WriteString(column)
		sets.WriteString("=?")

		i++
	}
	fieldNames.PushBack(this.columnToField[this.pk])
	sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s=?", this.tableName, sets.String(), this.pk)
	return
}

//预处理占位符的数据
func (this *BaseDao) prepareValues(data interface{}, fieldNames list.List) []interface{} {
	values := make([]interface{}, len(this.columnToField))
	object := reflect.ValueOf(data).Elem()

	i := 0
	for e := fieldNames.Front(); e != nil; e = e.Next() {
		name := e.Value.(string)
		field := object.FieldByName(name)
		values[i] = this.fieldValue(field)
		i++
	}

	return values
}

//reflect.Value获取值
func (this *BaseDao) fieldValue(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Ptr:
		m := v.MethodByName("Format")
		rets := m.Call([]reflect.Value{reflect.ValueOf(timeFormate)})
		t := rets[0].String()
		return t
		//return this.valueToString(v)
	default:
		return nil
	}
}

//根据SQL查询多条记录
func (this *BaseDao) Find() (*list.List, error) {
	var sql = "SELECT * FROM userinfo"
	rows, err := Open().Query(sql)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	data := list.New()
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}
		obj := this.parseQuery(columns, values)
		data.PushBack(obj)
	}
	return data, err
}

//对一条查询结果进行封装
func (this *BaseDao) parseQuery(columns []string, values []interface{}) interface{} {

	obj := reflect.New(this.EntityType).Interface()
	typ := reflect.ValueOf(obj).Elem()

	for i, col := range values {
		if col != nil {
			//fmt.Println(i, col)
			name := this.columnToField[columns[i]]
			//fmt.Println(name)
			var field = typ.FieldByName(name)

			//fmt.Println(field)
			b, ok := col.([]byte)
			if ok {
				this.parseQueryColumn(field, string(b))
			} else {
				this.parseQueryColumn(field, col)
			}
			//fmt.Println(field.String())
		}
	}
	return obj
}

//单个属性赋值
func (this *BaseDao) parseQueryColumn(field reflect.Value, s interface{}) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(reflect.ValueOf(s).String())
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(reflect.ValueOf(s).String(), 10, 0)
		field.SetUint(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(reflect.ValueOf(s).String(), 10, 0)
		field.SetInt(v)
	case reflect.Float32:
		v, _ := strconv.ParseFloat(reflect.ValueOf(s).String(), 32)
		field.SetFloat(v)
	case reflect.Float64:
		v, _ := strconv.ParseFloat(reflect.ValueOf(s).String(), 64)
		field.SetFloat(v)
	case reflect.Ptr:
		values := new(bytes.Buffer)
		vs := reflect.ValueOf(s)
		m := vs.MethodByName("Format")
		rets := m.Call([]reflect.Value{reflect.ValueOf(timeFormate)})
		t := rets[0].String()
		values.WriteString(t)
		v, _ := time.Parse(timeFormate, values.String())
		field.Set(reflect.ValueOf(&v))
	default:

	}
}
