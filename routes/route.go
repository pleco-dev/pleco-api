package routes

import (
	"go-auth-app/modules/auth"
	"go-auth-app/modules/user"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	userRepo := user.NewRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(nil)
	authService := auth.NewService(authRepo, userService)
	authHandler := auth.NewHandler(authService)

	auth.SetupRoutes(router.Group("/"), authHandler)
	user.SetupRoutes(router.Group("/"), userHandler)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
