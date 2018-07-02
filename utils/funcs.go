package utils

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
)

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
