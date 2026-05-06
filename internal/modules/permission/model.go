package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID        uint `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string
}
