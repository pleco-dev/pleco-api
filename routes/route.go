package routes

import (
	"go-auth-app/controllers"
	"go-auth-app/middleware"
	"go-auth-app/repositories"
	"go-auth-app/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// ✅ Repository
	userRepo := &repositories.UserRepoDB{}

	// ✅ Service
	authService := &services.AuthService{
		UserRepo: userRepo,
	}

	userService := &services.UserService{
		UserRepo: userRepo,
	}

	// ✅ Controller (pakai service, bukan repo lagi)
	authController := controllers.AuthController{
		AuthService: authService,
	}

	userController := controllers.UserController{
		UserService: userService,
	}

	api := router.Group("/api")

	// ========================
	// 🔓 PUBLIC ROUTES
	// ========================
	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)
	api.POST("/refresh", middleware.RefreshToken)

	// ========================
	// 🔐 PROTECTED ROUTES
	// ========================
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/profile", userController.Profile)
	protected.POST("/logout", authController.Logout)

	// ========================
	// 👑 ADMIN ROUTES
	// ========================
	admin := protected.Group("/admin")
	admin.Use(middleware.AdminOnly())

	admin.GET("/users", userController.GetAllUsers)
	admin.DELETE("/users/:id", userController.DeleteUser)
}
