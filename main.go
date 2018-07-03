package main

import (
	"admin/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	app.Router(router)
	router.Run(fmt.Sprintf(":%d", 8000))

}
