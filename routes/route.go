package routes

import (
	"go-auth-app/modules/auth"
	"go-auth-app/modules/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	userModule := user.BuildModule()
	authModule := auth.BuildModule(userModule.Service)

	auth.SetupRoutes(router.Group("/"), authModule.Handler)
	user.SetupRoutes(router.Group("/"), userModule.Handler)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
