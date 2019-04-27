package core

import (
	"github.com/go-sql-driver/mysql"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
	"time"
)

type Xorm struct {
	db *xorm.Engine
}

var instanceXorm *Xorm

func GetXormAuto() *xorm.Engine {
	db, err := GetXorm()
	if err != nil {
		panic("数据库访问错误")
	}
	return db.db
}

func GetXorm() (*Xorm, error) {
	if instanceXorm == nil {
		config := GetInstanceConfig()
		mysqlConfig := mysql.NewConfig()
		mysqlConfig.User = config.Db.User
		mysqlConfig.DBName = config.Db.Name
		mysqlConfig.Passwd = config.Db.Password
		mysqlConfig.Params = config.Db.Params
		mysqlConfig.Net = "tcp"
		mysqlConfig.Addr = config.Db.Host + ":" + config.Db.Port

		engine, err := xorm.NewEngine("mysql", mysqlConfig.FormatDSN())
		if config.Db.Prefix != "" {
			tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, config.Db.Prefix)
			engine.SetTableMapper(tbMapper)
		}

		instanceXorm = &Xorm{db: engine}
		instanceXorm.db.SetMaxOpenConns(config.Db.MaxOpenConn)
		if config.Db.Debug {
			instanceXorm.db.ShowSQL(true)
		}
		return instanceXorm, err
	}
	return instanceXorm, nil
}

func (x *Xorm) Ping() error {
	session := x.db.NewSession()
	defer session.Close()
	return session.DB().Ping()
}

func init() {
	go func() {
		for {
			if instanceXorm != nil {
				err := instanceXorm.Ping()
				if err != nil {
					println(err.Error())
					instanceXorm.db.Close()
					instanceXorm = nil
				}
			}
			time.Sleep(time.Second)
		}
	}()
}
