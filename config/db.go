package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var database *gorm.DB
	var err error

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		GetEnv("DB_HOST", "db"),
		GetEnv("DB_USER", "postgres"),
		GetEnv("DB_PASSWORD", "password"),
		GetEnv("DB_NAME", "auth_db"),
		GetEnv("DB_PORT", "5432"),
		GetEnv("DB_SSLMODE", "disable"),
	)

	// retry mechanism
	for i := 0; i < 10; i++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Database connected")
			break
		}

		log.Println("⏳ Waiting for database...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("❌ DB connection failed after retries: %v", err)
	}

	// optional: tetap simpan global (biar backward compatible)
	DB = database

	return database
}
