package permission

import (
	"pleco-api/internal/cache"

	"gorm.io/gorm"
)

type Module struct {
	Repository Repository
	Service    *Service
}

func BuildModule(db *gorm.DB, stores ...cache.Store) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	if len(stores) > 0 {
		service.Cache = stores[0]
	}

	return &Module{
		Repository: repository,
		Service:    service,
	}
}
