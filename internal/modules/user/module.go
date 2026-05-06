package user

import (
	"pleco-api/internal/cache"
	"pleco-api/internal/modules/audit"
	tokenModule "pleco-api/internal/modules/token"

	"gorm.io/gorm"
)

type Module struct {
	Repository Repository
	Service    *Service
	Handler    *Handler
}

func BuildModule(db *gorm.DB, auditSvc *audit.Service, stores ...cache.Store) *Module {
	repository := NewRepository(db)
	refreshRepo := tokenModule.NewRefreshTokenRepository(db)
	service := NewService(db, repository, refreshRepo, auditSvc)
	if len(stores) > 0 {
		service.Cache = stores[0]
	}
	handler := NewHandler(service)

	return &Module{
		Repository: repository,
		Service:    service,
		Handler:    handler,
	}
}
