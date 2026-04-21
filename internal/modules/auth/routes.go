package auth

import (
	"go-api-starterkit/internal/middleware"
	"go-api-starterkit/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup, handler *AuthHandler, jwtService *services.JWTService) {
	auth := api.Group("/auth")
	loginLimiter := middleware.NewRateLimiter(5, time.Minute)
	registerLimiter := middleware.NewRateLimiter(5, time.Minute)
	passwordLimiter := middleware.NewRateLimiter(3, 5*time.Minute)
	refreshLimiter := middleware.NewRateLimiter(10, time.Minute)
	socialLimiter := middleware.NewRateLimiter(5, time.Minute)

	auth.POST("/register", registerLimiter.Middleware(), handler.Register)
	auth.POST("/login", loginLimiter.Middleware(), handler.Login)
	auth.POST("/refresh", refreshLimiter.Middleware(), handler.RefreshToken)
	auth.GET("/verify", handler.VerifyEmail)
	auth.GET("/resend-verification", passwordLimiter.Middleware(), handler.ResendVerification)
	auth.POST("/forgot-password", passwordLimiter.Middleware(), handler.ForgotPassword)
	auth.POST("/reset-password", passwordLimiter.Middleware(), handler.ResetPassword)
	auth.POST("/social-login", socialLimiter.Middleware(), handler.SocialLogin)

	protected := auth.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	protected.GET("/profile", handler.Profile)
	protected.GET("/sessions", handler.ListSessions)
	protected.POST("/logout", handler.Logout)
	protected.POST("/logout-all", handler.LogoutAll)
	protected.POST("/logout-others", handler.LogoutOtherSessions)
	protected.DELETE("/sessions/:id", handler.RevokeSession)
}
