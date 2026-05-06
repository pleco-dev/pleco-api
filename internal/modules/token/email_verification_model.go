package token

import (
	"time"

	"gorm.io/gorm"
)

type EmailVerificationToken struct {
	ID        uint           `json:"id"`
	UserID    uint           `json:"user_id"`
	Token     string         `json:"-"` // SHA-256 hex of the raw token emailed to the user (never store the raw token)
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
