package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"mygin_websrv/conf"
)

var mysqlClt *xorm.Engine

func init() {
	if nil == mysqlClt {
		var err error
		mysqlClt, err = xorm.NewEngine(conf.Db["mysql"].DriverName, conf.Db["mysql"].Dsn)
		if err != nil {
			log.Fatal(err)
		}
		mysqlClt.SetMaxIdleConns(conf.Db["mysql"].MaxIdle) //空闲连接
		mysqlClt.SetMaxOpenConns(conf.Db["mysql"].MaxOpen) //最大连接数
		mysqlClt.ShowSQL(conf.Db["mysql"].ShowSql)
		mysqlClt.ShowExecTime(conf.Db["mysql"].ShowExecTime)
	}
}
