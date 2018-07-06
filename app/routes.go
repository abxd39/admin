package app

import (
	"admin/app/controller"
	"admin/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//session
	r.Use(sessions.Sessions("mysession", utils.Store))

	new(controller.PublicController).Router(r)
	new(controller.TestController).Router(r)
	new(controller.WebsocketController).Router(r)
	new(controller.AdminController).Router(r)
	new(controller.ContextController).Router(r)
	new(controller.WebUserManageController).Router(r)
	new(controller.CurrencyController).Router(r)
	new(controller.TokenController).Router(r)
}
