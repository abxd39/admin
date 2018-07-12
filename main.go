package main

import (
	"admin/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	app.Router(router)
	router.Run(fmt.Sprintf(":%d", 8001))

}

func ResetController() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
