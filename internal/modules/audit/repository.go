package audit

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Filter struct {
	Page        int
	Limit       int
	Action      string
	Resource    string
	Status      string
	ActorUserID *uint
	Search      string
	DateFrom    *time.Time
	DateTo      *time.Time
}

type Repository interface {
	Create(log *AuditLog) error
	FindAllWithFilter(filter Filter) ([]AuditLog, int64, error)
	FindForExport(filter Filter) ([]AuditLog, error)
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

func (r *gormRepository) FindAllWithFilter(filter Filter) ([]AuditLog, int64, error) {
	var (
		logs  []AuditLog
		total int64
	)

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}

	query := r.applyFilter(r.db.Model(&AuditLog{}), filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.Order("created_at DESC").Limit(filter.Limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *gormRepository) FindForExport(filter Filter) ([]AuditLog, error) {
	var logs []AuditLog

	query := r.applyFilter(r.db.Model(&AuditLog{}), filter)
	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *gormRepository) applyFilter(query *gorm.DB, filter Filter) *gorm.DB {
	if filter.Action != "" {
		query = query.Where("action = ?", filter.Action)
	}
	if filter.Resource != "" {
		query = query.Where("resource = ?", filter.Resource)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.ActorUserID != nil {
		query = query.Where("actor_user_id = ?", *filter.ActorUserID)
	}
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	if search := strings.TrimSpace(filter.Search); search != "" {
		like := "%" + search + "%"
		query = query.Where(
			"(action ILIKE ? OR resource ILIKE ? OR status ILIKE ? OR description ILIKE ? OR ip_address ILIKE ? OR user_agent ILIKE ?)",
			like, like, like, like, like, like,
		)
	}
	return query
}
