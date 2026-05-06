package role

import "gorm.io/gorm"

type Repository interface {
	FindByID(id uint) (*Role, error)
	FindAll() ([]Role, error)
	WithTx(tx *gorm.DB) Repository
}

type GormRepository struct {
	db *gorm.DB
}

var _ Repository = (*GormRepository)(nil)

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindByID(id uint) (*Role, error) {
	var role Role
	if err := r.db.Preload("RolePermissions").First(&role, id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *GormRepository) FindAll() ([]Role, error) {
	var roles []Role
	if err := r.db.Preload("RolePermissions").Order("id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *GormRepository) WithTx(tx *gorm.DB) Repository {
	return &GormRepository{db: tx}
}
