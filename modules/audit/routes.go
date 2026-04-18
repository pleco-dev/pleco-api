package audit

import (
	"go-auth-app/middleware"
	"go-auth-app/modules/permission"
	"go-auth-app/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup, handler *Handler, jwtService *services.JWTService, permissionService *permission.Service) {
	protected := api.Group("/auth")
	protected.Use(middleware.AuthMiddleware(jwtService))

	admin := protected.Group("/admin")
	admin.GET("/audit-logs", middleware.RequirePermission(permissionService, "audit.read"), handler.GetLogs)
}
