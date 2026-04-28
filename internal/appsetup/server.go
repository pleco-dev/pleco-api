package appsetup

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"pleco-api/internal/config"
	"pleco-api/internal/middleware"
	"pleco-api/internal/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func RunAPI(registerDocs func(*gin.Engine)) error {
	config.LoadEnv()

	appConfig := config.LoadAppConfig()
	if err := appConfig.Validate(); err != nil {
		return err
	}

	db := config.ConnectDB(appConfig.DatabaseURL)
	RunStartupTasks(appConfig, db)

	jwtService := services.NewJWTService(appConfig.JWTSecret)
	rateStore := newRateLimitStore(appConfig.RedisURL)
	router, err := BuildRouter(db, appConfig, jwtService, rateStore)
	if err != nil {
		return err
	}

	if registerDocs != nil {
		registerDocs(router)
	}

	srv := &http.Server{
		Addr:              ":" + appConfig.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on :%s", appConfig.Port)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case err := <-serverErr:
		return err
	case <-shutdownCtx.Done():
		log.Println("Shutdown signal received (SIGINT/SIGTERM)")
		log.Println("Gracefully shutting down HTTP server (max 10 seconds)...")
	}

	gracefulCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(gracefulCtx); err != nil {
		return err
	}

	if sqlDB, err := db.DB(); err == nil {
		if err := sqlDB.Close(); err != nil {
			return err
		}
	}

	if err := closeRateLimitStore(rateStore); err != nil {
		log.Printf("Error closing rate limit store: %v", err)
	}

	return <-serverErr
}

type rateLimitStoreCloser interface {
	middleware.RateLimitStore
	Close() error
}

func closeRateLimitStore(store middleware.RateLimitStore) error {
	closer, ok := store.(rateLimitStoreCloser)
	if !ok {
		return nil
	}
	return closer.Close()
}
