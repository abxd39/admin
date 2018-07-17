package main

import (
	"fmt"

	"admin/app"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// custom middleware
	//r.Use(utils.CheckLogin())

	app.Router(r)
	r.Run(fmt.Sprintf(":%d", 8001))
}
