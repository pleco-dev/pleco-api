package repositories

import (
	"go-auth-app/config"
	"go-auth-app/models"
)

// EmailVerificationTokenRepoDB implements EmailVerificationTokenRepository using a global DB instance.
type EmailVerificationTokenRepoDB struct{}

func NewEmailVerificationTokenRepo() EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepoDB{}
}

func (r *EmailVerificationTokenRepoDB) Create(token *models.EmailVerificationToken) error {
	return config.DB.Create(token).Error
}

func (r *EmailVerificationTokenRepoDB) FindByToken(token string) (*models.EmailVerificationToken, error) {
	var verification models.EmailVerificationToken
	err := config.DB.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *EmailVerificationTokenRepoDB) DeleteByID(id uint) error {
	return config.DB.Delete(&models.EmailVerificationToken{}, id).Error
}

func (r *EmailVerificationTokenRepoDB) DeleteByUserID(userID uint) error {
	return config.DB.Where("user_id = ?", userID).
		Delete(&models.EmailVerificationToken{}).Error
}
