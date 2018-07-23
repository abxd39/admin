package main

import (
	"fmt"
	"os"

	"admin/app"
	"admin/middleware"
	"admin/session"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("ADMIN_API_ENV") == "" {
		panic("环境变量ADMIN_API_ENV未设置")
	}

	r := gin.Default()

	// session
	r.Use(sessions.Sessions("mysession", session.Store))

	// custom middleware
	r.Use(middleware.JsCors())
	r.Use(middleware.CheckLogin())

	app.Router(r)
	r.Run(fmt.Sprintf(":%d", utils.Cfg.MustInt("http", "port")))
}
