package core

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

var instanceDb *gorm.DB

func GetDb() (*gorm.DB, error) {
	if instanceDb == nil {
		config := GetInstanceConfig()
		mysqlConfig := mysql.NewConfig()
		mysqlConfig.User = config.Db.User
		mysqlConfig.DBName = config.Db.Name
		mysqlConfig.Passwd = config.Db.Password
		mysqlConfig.Params = config.Db.Params
		mysqlConfig.Net = "tcp"
		mysqlConfig.Addr = config.Db.Host + ":" + config.Db.Port

		db, err := gorm.Open("mysql", mysqlConfig.FormatDSN())
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			if len(defaultTableName) > len(config.Db.Prefix) && defaultTableName[:len(config.Db.Prefix)] == config.Db.Prefix {
				return defaultTableName
			}
			return config.Db.Prefix + defaultTableName
		}
		db.DB().SetMaxOpenConns(config.Db.MaxOpenConn)
		db.SingularTable(true)
		instanceDb = db
		return instanceDb, err
	}
	return instanceDb, nil
}

func GetGormAuto() *gorm.DB {
	return GetDbAuto()
}

func GetDbAuto() *gorm.DB {
	db, err := GetDb()
	if err != nil {
		panic("数据库访问错误")
	}
	return db
}

func DbQuery(query string, args ...interface{}) []map[string]interface{} {
	db := GetDbAuto()
	rows, _ := db.DB().Query(query, args...) // Note: Ignoring errors for brevity
	defer rows.Close()
	cols, _ := rows.Columns()
	var data []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		data = append(data, m)
	}
	return data
}

var isUpdateComment = make(map[string]bool)

func GetDbAutoMigrate(values ...interface{}) *gorm.DB {
	db := GetDbAuto()
	db.AutoMigrate(values...)

	for _, value := range values {
		scope := db.NewScope(value)
		tableName := scope.TableName()
		_, isOk := isUpdateComment[tableName]
		if !isOk {
			field := scope.Fields()
			for e := range field {
				comment := field[e].Tag.Get("comment")

				if len(strings.Trim(comment, " ")) > 0 {
					fieldType := db.Dialect().DataTypeOf(field[e].StructField)
					scope.Raw(fmt.Sprintf("ALTER TABLE `%v` MODIFY COLUMN `%v` %v COMMENT '%v';", tableName, field[e].DBName, fieldType, comment)).Exec()
				}
			}
		}
		isUpdateComment[tableName] = true
	}

	return db
}

func init() {
	go func() {
		for {
			if instanceDb != nil {
				err := instanceDb.DB().Ping()
				if err != nil {
					println(err.Error())
					instanceDb.Close()
					instanceDb = nil
				}
			}
			time.Sleep(time.Second)
		}
	}()
}
