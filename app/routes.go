package app

import (
	"admin/app/controller"
	"admin/session"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//session
	r.Use(sessions.Sessions("mysession", session.Store))

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
	new(controller.ConfigController).Router(r)
}
