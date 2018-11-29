package main

import (
	"os"
	"admin/app"
	"admin/utils"
	"admin/cron"
	"admin/app/models"
	"fmt"
	"runtime"
	"sync"
)

var(
	counter int
	wg sync.WaitGroup
	mutex sync.Mutex
)

func main() {
	if os.Getenv("ADMIN_API_ENV") == ""{
		panic("环境变量ADMIN_API_ENV未设置")
	}
	wg.Add(2)
	go incCounter(1)
	go incCounter(2)
	wg.Wait()
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

func incCounter(id int){
	{
		fmt.Println("多加任意多个大括号都是可以的 " )
	}
	defer  wg.Done()
	for count:=0;count<2;count++{
		mutex.Lock()
		{
			value:=counter
			runtime.Gosched()
			value++
			counter=value
		}
		mutex.Unlock()
		{
			fmt.Println("这是什么情况啊啊")
		}
	}
}