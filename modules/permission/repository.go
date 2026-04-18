package permission

import "gorm.io/gorm"

type Repository interface {
	HasRolePermission(roleName, permission string) (bool, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) HasRolePermission(roleName, permission string) (bool, error) {
	var count int64

	err := r.db.Table("role_permissions").
		Joins("JOIN roles ON roles.id = role_permissions.role_id").
		Where("roles.name = ? AND role_permissions.permission = ?", roleName, permission).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
