package repositories

import "go-auth-app/models"

type PasswordResetRepository interface {
	Create(token *models.PasswordResetToken) error
	FindByToken(token string) (*models.PasswordResetToken, error)
	Delete(token string) error
}
