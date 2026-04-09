package token

import "go-auth-app/config"

type EmailVerificationRepository interface {
	Create(token *EmailVerificationToken) error
	FindByToken(token string) (*EmailVerificationToken, error)
	DeleteByID(id uint) error
	DeleteByUserID(userID uint) error
}

type GormEmailVerificationRepository struct{}

var _ EmailVerificationRepository = (*GormEmailVerificationRepository)(nil)

func NewEmailVerificationRepository() EmailVerificationRepository {
	return &GormEmailVerificationRepository{}
}

func (r *GormEmailVerificationRepository) Create(token *EmailVerificationToken) error {
	return config.DB.Create(token).Error
}

func (r *GormEmailVerificationRepository) FindByToken(token string) (*EmailVerificationToken, error) {
	var verification EmailVerificationToken
	if err := config.DB.Where("token = ?", token).First(&verification).Error; err != nil {
		return nil, err
	}

	return &verification, nil
}

func (r *GormEmailVerificationRepository) DeleteByID(id uint) error {
	return config.DB.Delete(&EmailVerificationToken{}, id).Error
}

func (r *GormEmailVerificationRepository) DeleteByUserID(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&EmailVerificationToken{}).Error
}
