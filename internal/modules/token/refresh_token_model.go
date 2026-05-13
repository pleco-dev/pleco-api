package token

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID             uint
	TokenHash          string
	FamilyID           string
	RotatedFromTokenID *uint
	ReplacedByTokenID  *uint
	DeviceID           string
	UserAgent          string
	IPAddress          string
	ExpiredAt          time.Time
	RevokedAt          *time.Time
	RevokeReason       string
}
