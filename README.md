# orm-engine

### using reflect do

1. orm insert
2. orm find
3. orm update

run result is in diretory `images`

### userinfo

```
type UserInfo struct {
	UID        int        `table:"userinfo" column:"uid"`
	UserName   string     `column:"username"`
	DepartName string     `column:"departname"`
	CreateAt   *time.Time `column:"created"`
}
```

### dao
```
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
```

### sql auto-mapping

insert auto-mapping
```
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
```
### parse from reflect value

```
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
```

### more about to see source code please!

### run mysql in docker
```
sysuygm@sysuygm:~/golang-workspace/src/orm-engine$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
agenda              latest              83b913da3c5f        9 days ago          762MB
golang              1.8                 ba52c9ef0f5c        13 days ago         712MB
mysql               5.7                 7d83a47ab2d2        2 weeks ago         408MB
sysuygm@sysuygm:~/golang-workspace/src/orm-engine$ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
cd42f3c3151a        mysql:5.7           "docker-entrypoint..."   36 hours ago        Up 12 hours         0.0.0.0:3306->3306/tcp   mysql2
sysuygm@sysuygm:~/golang-workspace/src/orm-engine$ 

```
### mysql database

create database and table

```
mysql> describe userinfo;
+------------+-------------+------+-----+---------+----------------+
| Field      | Type        | Null | Key | Default | Extra          |
+------------+-------------+------+-----+---------+----------------+
| uid        | int(10)     | NO   | PRI | NULL    | auto_increment |
| username   | varchar(64) | YES  |     | NULL    |                |
| departname | varchar(64) | YES  |     | NULL    |                |
| created    | datetime    | YES  |     | NULL    |                |
+------------+-------------+------+-----+---------+----------------+
4 rows in set (0.00 sec)
```

verify result

```
mysql> select * from userinfo;
+-----+----------+------------+---------------------+
| uid | username | departname | created             |
+-----+----------+------------+---------------------+
|   1 | sysu0dd  | depart0001 | 2017-12-26 11:36:03 |
+-----+----------+------------+---------------------+
1 row in set (0.00 sec)

mysql> select * from userinfo;
+-----+----------+------------+---------------------+
| uid | username | departname | created             |
+-----+----------+------------+---------------------+
|   1 | sysu0dd  | depart0001 | 2017-12-26 11:36:03 |
|   2 | sysu0dd  | depart0001 | 2017-12-26 11:38:37 |
|   3 | sysu0dd  | depart0001 | 2017-12-26 11:39:47 |
|   4 | sysu0dd  | depart0001 | 2017-12-26 11:47:34 |
|   5 | sysu0dd  | depart0001 | 2017-12-26 11:48:04 |
|   6 | sysu0dd  | depart0001 | 2017-12-26 11:49:02 |
|   7 | sysu0dd  | depart0001 | 2017-12-26 11:50:16 |
|   8 | sysu0dd  | depart0001 | 2017-12-26 11:50:44 |
|   9 | sysu0dd  | depart0001 | 2017-12-26 11:50:59 |
|  10 | sysu0dd  | depart0001 | 2017-12-26 11:51:17 |
|  11 | sysu0dd  | depart0001 | 2017-12-26 11:51:59 |
|  12 | sysu0dd  | depart0001 | 2017-12-26 11:54:47 |
|  13 | sysu0dd  | depart0001 | 2017-12-26 11:56:12 |
|  14 | sysu0dd  | depart0001 | 2017-12-26 11:56:56 |
|  15 | sysu0dd  | depart0001 | 2017-12-26 11:58:04 |
|  16 | sysu0dd  | depart0001 | 2017-12-26 11:58:58 |
|  17 | sysu0dd  | depart0001 | 2017-12-26 11:59:32 |
|  18 | sysu0dd  | depart0001 | 2017-12-26 11:59:55 |
+-----+----------+------------+---------------------+
18 rows in set (0.01 sec)

mysql> 
```

### run 

result is ok!

```
sysuygm@sysuygm:~/golang-workspace/src/orm-engine$ go run main.go 
INSERT INTO userinfo (created,uid,username,departname) VALUES (?,?,?,?)   [2017-12-26 11:59:55 0 sysu0dd depart0001]
result: &{1 sysu0dd depart0001 2017-12-26 11:36:03 +0000 UTC}
result: &{2 sysu0dd depart0001 2017-12-26 11:38:37 +0000 UTC}
result: &{3 sysu0dd depart0001 2017-12-26 11:39:47 +0000 UTC}
result: &{4 sysu0dd depart0001 2017-12-26 11:47:34 +0000 UTC}
result: &{5 sysu0dd depart0001 2017-12-26 11:48:04 +0000 UTC}
result: &{6 sysu0dd depart0001 2017-12-26 11:49:02 +0000 UTC}
result: &{7 sysu0dd depart0001 2017-12-26 11:50:16 +0000 UTC}
result: &{8 sysu0dd depart0001 2017-12-26 11:50:44 +0000 UTC}
result: &{9 sysu0dd depart0001 2017-12-26 11:50:59 +0000 UTC}

```
