package main

import (
	"admin/app"
	"admin/app/models"
	"admin/cron"
	"admin/utils"
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"admin/session"
	"admin/middleware"

)

func main() {
	if os.Getenv("ADMIN_API_ENV") == "" {
		panic("环境变量ADMIN_API_ENV未设置")
	}

	// 定时任务
	cron.InitCron()
	go models.DailyStart()

	//启动定时器
	//go new(models.TokenFeeDailySheet).BoottimeTimingSettlement()
	//go new(models.WalletInoutDailySheet).BoottimeTimingSettlement()
	//go new(models.CurencyFeeDailySheet).BoottimeTimingSettlement()

	// 配置gin
	r := gin.Default()
	//定时任务
	go models.DailyStart()
	//定时任务工具
	//go models.DailyStart1()
	// session
	r.Use(sessions.Sessions("mysession", session.Store))

	// custom middleware
	r.Use(middleware.JsCors())
	r.Use(middleware.CheckLogin())
	app.Router(r)

	r.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))

}

//func main() {
//
//	models.DailyStart()
//
//}
