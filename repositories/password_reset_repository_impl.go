package repositories

import (
	"go-auth-app/config"
	"go-auth-app/models"
)

// PasswordResetRepoDB implements PasswordResetRepository using a global DB instance.
type PasswordResetRepoDB struct{}

func NewPasswordResetRepo() PasswordResetRepository {
	return &PasswordResetRepoDB{}
}

func (r *PasswordResetRepoDB) Create(token *models.PasswordResetToken) error {
	return config.DB.Create(token).Error
}

func (r *PasswordResetRepoDB) FindByToken(token string) (*models.PasswordResetToken, error) {
	var reset models.PasswordResetToken
	err := config.DB.Where("token = ?", token).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *PasswordResetRepoDB) Delete(token string) error {
	return config.DB.Where("token = ?", token).
		Delete(&models.PasswordResetToken{}).Error
}
