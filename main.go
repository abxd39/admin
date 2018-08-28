package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"admin/app"
	"admin/utils"
	"admin/middleware"
	"admin/cron"
	"admin/app/models"
	"admin/session"
	"fmt"
)

func main() {
	if os.Getenv("ADMIN_API_ENV") == "" {
		panic("环境变量ADMIN_API_ENV未设置")
	}


	// 定时任务
	cron.InitCron()
	go models.DailyStart()

	// 配置gin
	r := gin.Default()
	// session
	r.Use(sessions.Sessions("mysession", session.Store))

	// custom middleware
	r.Use(middleware.JsCors())
	r.Use(middleware.CheckLogin())

	// 启动gin
	app.Init()
	app.App.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))

}

//func main() {
//
//	models.DailyStart1()
//
//}
