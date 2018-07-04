package utils

import (
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var Engine_currency *xorm.Engine
var Engine_token *xorm.Engine
var Engine_common *xorm.Engine
var Engine_context *xorm.Engine
var Engine_backstage *xorm.Engine
var Redis *redis.Conn

func init() {
	var err error

	//mysql初始化
	dsource := "ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_currency?charset=utf8"
	Engine_currency, err = xorm.NewEngine("mysql", dsource)
	if err != nil {
		panic(err)
	}
	err = Engine_currency.Ping()
	if err != nil {
		panic(err)
	}

	dsource = "root:current@tcp(47.106.136.96:3306)/g_token?charset=utf8"
	Engine_token, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	dsource = "root:current@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	Engine_common, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_common.Ping()
	if err != nil {
		panic(err)
	}
	//context manage
	dsource = "ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	//dsource = "conn=ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	Engine_context, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_context.Ping()
	if err != nil {
		panic(err)
	}

	dsource = "ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_backstage?charset=utf8"
	//dsource = "conn=ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	Engine_backstage, err = xorm.NewEngine("mysql", dsource)

	if err != nil {
		panic(err)
	}
	err = Engine_backstage.Ping()
	if err != nil {
		panic(err)
	}
	//redis初始化
	Redis = nil

}
