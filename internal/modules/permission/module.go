package permission

import "gorm.io/gorm"

type Module struct {
	Repository Repository
	Service    *Service
}

func BuildModule(db *gorm.DB) *Module {
	repository := NewRepository(db)
	service := NewService(repository)

	return &Module{
		Repository: repository,
		Service:    service,
	}
}
