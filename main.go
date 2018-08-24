package main

import (
	"admin/app"
	"admin/app/models"
	"admin/cron"
	"admin/middleware"
	"admin/session"
	"admin/utils"
	"fmt"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	//定时任务
	go models.DailyStart()
	//定时任务工具
	//go models.DailyStart1()
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
//	models.DailyStart()
//
//}
