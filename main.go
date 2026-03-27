package main

import (
	"go-auth-app/config"
	"go-auth-app/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.InitJWT()

	// Connect DB
	config.ConnectDB()

	// Seed admin
	config.SeedAdmin()

	// Init router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router)

	// Run server
	router.Run(":8080")
}
