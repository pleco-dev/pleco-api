package config

import (
	"context"
	"database/sql"
	"log"
	"net"
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

	// 🚀 Force IPv4 by overriding default dialer
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	net.DefaultResolver = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.DialContext(ctx, "tcp4", address)
		},
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
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db
}
