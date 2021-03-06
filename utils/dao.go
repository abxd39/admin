package utils

import (
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"admin/apis"
)

var Engine_currency *xorm.Engine
var Engine_token *xorm.Engine
var Engine_common *xorm.Engine
var Engine_context *xorm.Engine
var Engine_backstage *xorm.Engine
var Engine_wallet *xorm.Engine
var Redis *redis.Client
var AliClient *oss.Client

func init() {
	var err error

	//mysql初始化
	//dsource := "ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_currency?charset=utf8"
	Engine_currency, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "currency"))
	if err != nil {
		panic(err)
	}
	Engine_currency.ShowSQL(true)
	err = Engine_currency.Ping()
	if err != nil {
		panic(err)
	}

	//dsource = "root:current@tcp(47.106.136.96:3306)/g_token?charset=utf8"
	Engine_token, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "token"))

	if err != nil {
		panic(err)
	}
	Engine_token.ShowSQL(true)
	err = Engine_token.Ping()
	if err != nil {
		panic(err)
	}

	//dsource = "root:current@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	Engine_common, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "common"))
	if err != nil {
		panic(err)
	}
	Engine_common.ShowSQL(true)
	err = Engine_common.Ping()
	if err != nil {
		panic(err)
	}
	//context manage
	//dsource = "ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	//dsource = "conn=ccbk:ecrf981@@tcp(47.106.136.96:3306)/g_common?charset=utf8"
	Engine_context, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "context"))

	if err != nil {
		panic(err)
	}
	Engine_context.ShowSQL(true)
	err = Engine_context.Ping()
	if err != nil {
		panic(err)
	}

	Engine_backstage, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "backstage"))
	if err != nil {
		panic(err)
	}
	Engine_backstage.ShowSQL(true)
	err = Engine_backstage.Ping()
	if err != nil {
		panic(err)
	}
	//dsource := Cfg.MustValue("mysql", "wallet_conn")
	Engine_wallet, err = xorm.NewEngine("mysql", Cfg.MustValue("mysql", "wallet"))
	if err != nil {
		panic(err)
	}
	Engine_wallet.ShowSQL(true)
	err = Engine_wallet.Ping()
	if err != nil {
		panic(err)
	}
	//redis初始化
	client := redis.NewClient(&redis.Options{
		Addr:     Cfg.MustValue("redis", "addr"),
		Password: Cfg.MustValue("redis", "pwd"),
		DB:       0, // use default DB
	})

	_, err = client.Ping().Result()
	if err != nil {
		fmt.Printf(err.Error())
		panic(err)
	}
	Redis = client

	//ali
	AliClient, err = oss.New("http://oss-cn-shenzhen.aliyuncs.com", "LTAIcJgRedhxruPq", "d7p6tWRfy0B2QaRXk7q4mb5seLROtb")
	if err != nil {
		// HandleError(err)
		panic(err)
	}
	remote:= Cfg.MustValue("userUrl", "remoteUrl")
	local:=Cfg.MustValue("userUrl", "localUrl")
	key :=Cfg.MustValue("userUrl", "key")
	if remote==``{
		fmt.Printf(err.Error())
		panic(err)
	}
	apis.InitUserUrl(remote,local,key)
	award:=Cfg.MustValue("award","url")
	awardKey:=Cfg.MustValue("award","key")
	verifyKey:=Cfg.MustValue("award","verifykey")
	if award ==``{
		fmt.Printf(err.Error())
		panic(err)
	}
	apis.InitAwardUrl(award,awardKey,verifyKey)
	return

}
