package main

import (
	"fmt"

	"admin/app"
	"admin/middleware"
	"admin/session"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// session
	r.Use(sessions.Sessions("mysession", session.Store))

	// custom middleware
	r.Use(middleware.CheckLogin())

	app.Router(r)
	r.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))
}
