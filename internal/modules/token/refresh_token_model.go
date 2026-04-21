package token

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint
	TokenHash string
	DeviceID  string
	UserAgent string
	IPAddress string
	ExpiredAt time.Time
}
