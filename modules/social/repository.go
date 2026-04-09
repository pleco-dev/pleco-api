package social

import (
	"errors"

	"go-auth-app/config"

	"gorm.io/gorm"
)

type Repository interface {
	Create(socialAccount *SocialAccount) error
	FindByProvider(provider string, providerID string) (*SocialAccount, error)
}

type GormRepository struct{}

var _ Repository = (*GormRepository)(nil)

func NewRepository() Repository {
	return &GormRepository{}
}

func (r *GormRepository) Create(socialAccount *SocialAccount) error {
	if socialAccount == nil {
		return errors.New("socialAccount cannot be nil")
	}

	return config.DB.Create(socialAccount).Error
}

func (r *GormRepository) FindByProvider(provider, providerUserID string) (*SocialAccount, error) {
	if provider == "" || providerUserID == "" {
		return nil, errors.New("provider and providerUserID cannot be empty")
	}

	var account SocialAccount
	err := config.DB.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}
