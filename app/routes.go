package app

import (
	"admin/app/controller"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	new(controller.PublicController).Router(r)
	new(controller.TestController).Router(r)
	new(controller.WebsocketController).Router(r)
	new(controller.AdminController).Router(r)
	new(controller.ContextController).Router(r)
}

