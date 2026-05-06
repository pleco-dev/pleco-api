package social

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	Create(socialAccount *SocialAccount) error
	FindByProvider(provider string, providerID string) (*SocialAccount, error)
	FindByUserAndProvider(userID uint, provider string) (*SocialAccount, error)
	WithTx(tx *gorm.DB) Repository
}

type GormRepository struct {
	db *gorm.DB
}

var _ Repository = (*GormRepository)(nil)

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(socialAccount *SocialAccount) error {
	if socialAccount == nil {
		return errors.New("socialAccount cannot be nil")
	}

	return r.db.Create(socialAccount).Error
}

func (r *GormRepository) FindByProvider(provider, providerUserID string) (*SocialAccount, error) {
	if provider == "" || providerUserID == "" {
		return nil, errors.New("provider and providerUserID cannot be empty")
	}

	var account SocialAccount
	err := r.db.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *GormRepository) FindByUserAndProvider(userID uint, provider string) (*SocialAccount, error) {
	if userID == 0 || provider == "" {
		return nil, errors.New("userID and provider cannot be empty")
	}

	var account SocialAccount
	err := r.db.Where("user_id = ? AND provider = ?", userID, provider).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *GormRepository) WithTx(tx *gorm.DB) Repository {
	return &GormRepository{db: tx}
}
