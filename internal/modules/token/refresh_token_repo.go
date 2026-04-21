package token

import "gorm.io/gorm"

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	FindByID(id uint) (*RefreshToken, error)
	FindByUserAndDevice(userID uint, deviceID string) (*RefreshToken, error)
	FindByUser(userID uint) ([]RefreshToken, error)
	DeleteByID(id uint) error
	DeleteByUserAndID(userID, id uint) error
	DeleteByUser(userID uint) error
	DeleteByUserExceptDevice(userID uint, deviceID string) error
}

type GormRefreshTokenRepository struct {
	db *gorm.DB
}

var _ RefreshTokenRepository = (*GormRefreshTokenRepository)(nil)

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &GormRefreshTokenRepository{db: db}
}

func (r *GormRefreshTokenRepository) Save(token *RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *GormRefreshTokenRepository) FindByID(id uint) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.First(&token, id).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *GormRefreshTokenRepository) FindByUserAndDevice(userID uint, deviceID string) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *GormRefreshTokenRepository) FindByUser(userID uint) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens).Error
	return tokens, err
}

func (r *GormRefreshTokenRepository) DeleteByID(id uint) error {
	return r.db.Delete(&RefreshToken{}, id).Error
}

func (r *GormRefreshTokenRepository) DeleteByUserAndID(userID, id uint) error {
	return r.db.Where("user_id = ? AND id = ?", userID, id).Delete(&RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) DeleteByUserExceptDevice(userID uint, deviceID string) error {
	query := r.db.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id <> ?", deviceID)
	}
	return query.Delete(&RefreshToken{}).Error
}
