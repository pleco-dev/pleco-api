package auth

import (
	"go-api-starterkit/internal/config"
	"go-api-starterkit/internal/modules/audit"
	userModule "go-api-starterkit/internal/modules/user"
	"go-api-starterkit/internal/services"

	"gorm.io/gorm"
)

type Module struct {
	Service AuthService
	Handler *AuthHandler
}

func BuildModule(db *gorm.DB, cfg config.AppConfig, userService *userModule.Service, jwtService *services.JWTService, auditService *audit.Service) *Module {
	service := NewService(db, cfg, userService, jwtService, auditService)
	handler := NewHandler(service)

	return &Module{
		Service: service,
		Handler: handler,
	}
}
