package main

import (
	"admin/app"
	"admin/app/models"
	"admin/cron"
	"admin/utils"
	"fmt"
	"os"
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

	// 启动gin
	app.Init()
	app.App.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))
}

//func main() {
//
//	models.DailyStart()
//
//}
