package main

import (
	"os"
	"admin/app"
	"admin/utils"
	"admin/cron"
	"admin/app/models"
	"fmt"
)

func main() {
	if os.Getenv("ADMIN_API_ENV") == "" {
		panic("环境变量ADMIN_API_ENV未设置")
	}


	// 定时任务
	cron.InitCron()
	go models.DailyStart()


	// 启动gin
	app.Init()
	app.App.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))

}

//func main() {
//
//	models.DailyStart1()
//
//}
