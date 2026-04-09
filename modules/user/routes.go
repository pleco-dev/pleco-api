package user

import (
	"go-auth-app/config"
	"go-auth-app/middleware"
	"go-auth-app/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup, handler *Handler) {
	jwtService := services.NewJWTService(config.JWTSecret)

	protected := api.Group("/auth")
	protected.Use(middleware.AuthMiddleware(jwtService))

	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())
	admin.GET("/users", handler.GetAllUsers)
	admin.DELETE("/users/:id", handler.DeleteUser)
}
