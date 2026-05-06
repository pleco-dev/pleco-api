package role

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID              uint `json:"id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Name            string
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
}

type RolePermission struct {
	ID         uint `json:"id"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	RoleID     uint           `json:"role_id"`
	Permission string         `json:"permission"`
}
