package app

import (
	"admin/app/controller"
	"admin/middleware"
	"admin/session"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var App *gin.Engine

func Init() {
	App = gin.Default()

	// session
	App.Use(sessions.Sessions("mysession", session.Store))

	// 自定义中间件
	App.Use(middleware.JsCors())
	App.Use(middleware.CheckLogin())

	// 路由
	router(App)
}

func router(r *gin.Engine) {
	new(controller.PublicController).Router(r)
	new(controller.TestController).Router(r)
	new(controller.WebsocketController).Router(r)
	new(controller.AdminController).Router(r)
	new(controller.ContextController).Router(r)
	new(controller.WebUserManageController).Router(r)
	new(controller.CurrencyController).Router(r)
	new(controller.TokenController).Router(r)
	new(controller.RoleController).Router(r)
	new(controller.NodeController).Router(r)
	new(controller.NodeAPIController).Router(r)
	new(controller.ConfigController).Router(r)
	new(controller.WallectController).Router(r)
}
