package config

import (
	"database/sql"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DatabaseURL() string {
	return GetEnv("DATABASE_URL", "")
}

func ConnectDB(dsn string) *gorm.DB {
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL is not set")
	}

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		sqlDB, err2 := sql.Open("pgx", dsn)
		if err2 != nil {
			log.Println("⏳ Opening DB failed:", err2)
			time.Sleep(2 * time.Second)
			continue
		}

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})

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

	// ✅ Connection pool config
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(envInt("DB_MAX_OPEN_CONNS", 5))
	sqlDB.SetMaxIdleConns(envInt("DB_MAX_IDLE_CONNS", 2))
	sqlDB.SetConnMaxLifetime(time.Duration(envInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute)

	return db
}
