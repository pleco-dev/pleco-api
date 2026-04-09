package role

import "go-auth-app/config"

type Repository interface {
	FindByID(id uint) (*Role, error)
}

type GormRepository struct{}

var _ Repository = (*GormRepository)(nil)

func NewRepository() Repository {
	return &GormRepository{}
}

func (r *GormRepository) FindByID(id uint) (*Role, error) {
	var role Role
	if err := config.DB.First(&role, id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
