package repositories

import "go-auth-app/models"

type EmailVerificationTokenRepository interface {
	Create(token *models.EmailVerificationToken) error
	FindByToken(token string) (*models.EmailVerificationToken, error)
	DeleteByID(id uint) error
	DeleteByUserID(userID uint) error
}
