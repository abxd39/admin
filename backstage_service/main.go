package main

import (
	cf "admin/backstage_service/conf"
	"admin/backstage_service/dao"
	"admin/backstage_service/log"
	"admin/backstage_service/rpc"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	cf.Init()
	log.InitLog()
	dao.InitDao()
	log.Log.Infof("begin run backstage servce")
	go rpc.RPCServerInit()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Log.Infof("backstage server close by sig %s", sig.String())
	fmt.Println("hgehhahah ")
}
