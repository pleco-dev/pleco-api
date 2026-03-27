package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-auth-app/models"
)

var DB *gorm.DB

func ConnectDB() {
	// loadDotEnv()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "go_auth"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	log.Println("Database connected")

	database.AutoMigrate(&models.User{})

	DB = database
}

// func loadDotEnv() {
// 	// First try the default behavior (current working directory).
// 	if err := godotenv.Load(); err == nil {
// 		return
// 	}

// 	// Then try loading from the project root (one directory above `config/`).
// 	_, thisFile, _, ok := runtime.Caller(0)
// 	if !ok {
// 		return
// 	}
// 	rootDir := filepath.Dir(filepath.Dir(thisFile))
// 	_ = godotenv.Load(filepath.Join(rootDir, ".env"))
// }

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
