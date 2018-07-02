package app

import (
	"admin/app/controller"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"admin/utils"
)

func Router(r *gin.Engine) {
	//session
	r.Use(sessions.Sessions("mysession", utils.Store))

	new(controller.PublicController).Router(r)
	new(controller.TestController).Router(r)
	new(controller.WebsocketController).Router(r)
	new(controller.AdminController).Router(r)
	new(controller.ContextController).Router(r)
}

