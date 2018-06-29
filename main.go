package main

import (
	"admin/app"
	ini "admin/conf"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	flag.Parse()
	ini.Init()
	router := gin.Default()
	app.Router(router)
	router.Run(fmt.Sprintf(":%d", 8000))

}
