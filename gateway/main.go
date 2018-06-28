package main

import (
	cf "admin/gateway/conf"
	"admin/gateway/http"
	"admin/gateway/log"
	"admin/gateway/rpc"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cf.Init()
	log.InitLog()
	go http.InitHttpServer()
	go rpc.InitInnerService()

	quitChan := make(chan os.Signal)
	signal.Notify(quitChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	)

	sig := <-quitChan
	log.Log.Infof("server close by sig %s", sig.String())
	fmt.Println("Hello this is gateway!!")
}
