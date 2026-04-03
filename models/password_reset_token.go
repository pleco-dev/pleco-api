package models

import "time"

type PasswordResetToken struct {
	ID        uint
	UserID    uint
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
