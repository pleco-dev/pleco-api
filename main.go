package main

import (
	"go-auth-app/config"
	"go-auth-app/modules/auth"
	"go-auth-app/modules/user"
	"os"

	"github.com/gin-gonic/gin"
)

func initApp() *gin.Engine {
	// Load env & init JWT
	config.LoadEnv()
	config.InitJWT()

	// Connect DB
	config.ConnectDB()

	router := gin.Default()
	api := router.Group("/")

	userModule := user.BuildModule()
	authModule := auth.BuildModule(userModule.Service)

	// ===== ROUTES =====
	auth.SetupRoutes(api, authModule.Handler)
	user.SetupRoutes(api, userModule.Handler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}

func main() {
	router := initApp()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
