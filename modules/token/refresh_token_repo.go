package token

import "go-auth-app/config"

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	FindByUserAndDevice(userID uint, deviceID string) (*RefreshToken, error)
	FindByUser(userID uint) ([]RefreshToken, error)
	DeleteByID(id uint) error
	DeleteByUser(userID uint) error
}

type GormRefreshTokenRepository struct{}

var _ RefreshTokenRepository = (*GormRefreshTokenRepository)(nil)

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &GormRefreshTokenRepository{}
}

func (r *GormRefreshTokenRepository) Save(token *RefreshToken) error {
	return config.DB.Create(token).Error
}

func (r *GormRefreshTokenRepository) FindByUserAndDevice(userID uint, deviceID string) (*RefreshToken, error) {
	var token RefreshToken
	if err := config.DB.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *GormRefreshTokenRepository) FindByUser(userID uint) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := config.DB.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

func (r *GormRefreshTokenRepository) DeleteByID(id uint) error {
	return config.DB.Delete(&RefreshToken{}, id).Error
}

func (r *GormRefreshTokenRepository) DeleteByUser(userID uint) error {
	return config.DB.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}
