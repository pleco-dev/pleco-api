package audit

import "gorm.io/gorm"

type Repository interface {
	Create(log *AuditLog) error
	FindAllWithFilter(page, limit int, action, resource string) ([]AuditLog, int64, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(log *AuditLog) error {
	return r.db.Create(log).Error
}

func (r *gormRepository) FindAllWithFilter(page, limit int, action, resource string) ([]AuditLog, int64, error) {
	var (
		logs  []AuditLog
		total int64
	)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	query := r.db.Model(&AuditLog{})
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}
