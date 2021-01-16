package onebot

import (
	"database/sql"
	"fmt"
	"reflect"
	"yaya/core"

	_ "github.com/mattn/go-sqlite3"
)

// runDB 创建各个bot对应的数据库
func (conf *Yaml) runDB() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[数据库] DB =X=> =X=> Start Error: %v", err)
		}
	}()
	for i, _ := range conf.BotConfs {
		conf.BotConfs[i].DBPath = AppPath + core.Int2Str(conf.BotConfs[i].Bot) + "/XQ.db"
		CreatePath(conf.BotConfs[i].DBPath)
		conf.BotConfs[i].dbCreate(&XEvent{})
	}
}

// dbCreate 根据结构体生成数据库table，tag为"id"为主键，自增
func (bot *BotYaml) dbCreate(objptr interface{}) {
	db, err := sql.Open("sqlite3", bot.DBPath)
	if err != nil {
		panic(err)
	}

	table := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (", struct2name(objptr))
	for i, column := range strcut2columns(objptr) {
		table += fmt.Sprintf(" %s %s NULL", column, column2type(objptr, column))
		if i+1 != len(strcut2columns(objptr)) {
			table += ","
		} else {
			table += " );"
		}
	}
	if _, err := db.Exec(table); err != nil {
		panic(err)
	}
	bot.DB = db
}

// dbInsert 根据结构体插入一条数据
func (bot *BotYaml) dbInsert(objptr interface{}) int64 {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[数据库] DB =X=> =X=> Insert Error: %v", err)
		}
	}()
	rows, err := bot.DB.Query("SELECT * FROM " + struct2name(objptr))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()

	index := -1
	names := "("
	insert := "("
	for i, column := range columns {
		if column == "id" {
			index = i
			continue
		}
		if i+1 != len(columns) {
			names += column + ","
			insert += "?,"
		} else {
			names += column + ")"
			insert += "?)"
		}
	}

	stmt, err := bot.DB.Prepare("INSERT INTO " + struct2name(objptr) + names + " values " + insert)
	if err != nil {
		panic(err)
	}

	value := []interface{}{}
	if index == -1 {
		value = append(value, struct2values(objptr, columns)...)
	} else {
		value = append(value, append(struct2values(objptr, columns)[:index], struct2values(objptr, columns)[index+1:]...)...)
	}
	res, err := stmt.Exec(value...)
	if err != nil {
		panic(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	return id
}

// dbSelect 根据结构体查询对应的表，cmd可为" id = 0 "
func (bot *BotYaml) dbSelect(objptr interface{}, cmd string) {
	rows, err := bot.DB.Query(fmt.Sprintf("SELECT * FROM %s where %s", struct2name(objptr), cmd))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		columns, err := rows.Columns()
		if err != nil {
			panic(err)
		}
		err = rows.Scan(struct2addrs(objptr, columns)...)
		if err != nil {
			panic(err)
		}
	}
}

// strcut2columns 反射得到结构体的 tag 数组
func strcut2columns(objptr interface{}) []string {
	var columns []string
	elem := reflect.ValueOf(objptr).Elem()
	for i, flen := 0, elem.Type().NumField(); i < flen; i++ {
		columns = append(columns, elem.Type().Field(i).Tag.Get("db"))
	}
	return columns
}

// struct2name 反射得到结构体的名字
func struct2name(objptr interface{}) string {
	return reflect.ValueOf(objptr).Elem().Type().Name()
}

// column2type 反射得到结构体对应 tag 的 数据库数据类型
func column2type(objptr interface{}, column string) string {
	type_ := ""
	elem := reflect.ValueOf(objptr).Elem()
	for i, flen := 0, elem.Type().NumField(); i < flen; i++ {
		if column == elem.Type().Field(i).Tag.Get("db") {
			type_ = elem.Field(i).Type().String()
		}
	}
	if column == "id" {
		return "INTEGER PRIMARY KEY"
	}
	switch type_ {
	case "int64":
		return "INT"
	case "string":
		return "TEXT"
	default:
		return "TEXT"
	}
}

// struct2addrs 反射得到结构体对应数据库字段的属性地址
func struct2addrs(objptr interface{}, columns []string) []interface{} {
	var addrs []interface{}
	elem := reflect.ValueOf(objptr).Elem()
	for _, column := range columns {
		for i, flen := 0, elem.Type().NumField(); i < flen; i++ {
			if column == elem.Type().Field(i).Tag.Get("db") {
				addrs = append(addrs, elem.Field(i).Addr().Interface())
			}
		}
	}
	return addrs
}

// struct2values 反射得到结构体对应数据库字段的属性值
func struct2values(objptr interface{}, columns []string) []interface{} {
	var values []interface{}
	elem := reflect.ValueOf(objptr).Elem()
	for _, column := range columns {
		for i, flen := 0, elem.Type().NumField(); i < flen; i++ {
			if column == elem.Type().Field(i).Tag.Get("db") {
				switch elem.Field(i).Type().String() {
				case "int64":
					values = append(values, elem.Field(i).Int())
				case "string":
					values = append(values, elem.Field(i).String())
				default:
					values = append(values, elem.Field(i).String())
				}
			}
		}
	}
	return values
}
