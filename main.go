package main

import (
	"go-auth-app/config"
	"go-auth-app/models"
	"go-auth-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitJWT()

	// Connect DB
	config.ConnectDB()

	config.DB.AutoMigrate(&models.RefreshToken{})

	config.DB.AutoMigrate(&models.EmailVerificationToken{})

	// Seed admin
	config.SeedAdmin()

	// Init router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Run server
	router.Run(":8080")
}
