package models

import "time"

type EmailVerificationToken struct {
	ID        uint
	UserID    uint
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
