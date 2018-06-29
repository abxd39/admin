package app

import (
	"admin/app/controller"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func Router(r *gin.Engine) {
	new(controller.PublicController).Router(r)
	new(controller.TestController).Router(r)
	new(controller.WebsocketController).Router(r)
	new(controller.AdminController).Router(r)
	new(controller.ContextController).Router(r)
}

func CheckLogin() gin.HandlerFunc {
	fmt.Println("this is middleware  !!")
	return func(c *gin.Context) {
		value, err := uuid.NewV4()
		if err != nil {
			//写日志
			log.Fatalln("failed")
		}
		c.Writer.Header().Set("X-Request-Id", value.String())
		c.Next()

	}
} //<code class="go hljs">
