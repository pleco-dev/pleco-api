package social

import (
	"time"

	"gorm.io/gorm"
)

type SocialAccount struct {
	ID             uint `json:"id"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserID         uint
	Provider       string
	ProviderUserID string
	AvatarURL      string
}
