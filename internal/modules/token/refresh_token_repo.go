package token

import (
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Save(token *RefreshToken) error
	FindByID(id uint) (*RefreshToken, error)
	FindByUserAndDevice(userID uint, deviceID string) (*RefreshToken, error)
	FindByTokenHash(tokenHash string) (*RefreshToken, error)
	FindByUser(userID uint) ([]RefreshToken, error)
	RevokeByID(id uint, replacedByTokenID *uint, reason string) error
	RevokeFamily(userID uint, familyID string, reason string) error
	DeleteByID(id uint) error
	DeleteByUserAndID(userID, id uint) error
	DeleteByUser(userID uint) error
	DeleteByUserAndDevice(userID uint, deviceID string) error
	DeleteByUserExceptDevice(userID uint, deviceID string) error
	WithTx(tx *gorm.DB) RefreshTokenRepository
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
	if err := r.db.Where("user_id = ? AND device_id = ? AND revoked_at IS NULL", userID, deviceID).Order("created_at DESC").First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *GormRefreshTokenRepository) FindByTokenHash(tokenHash string) (*RefreshToken, error) {
	var token RefreshToken
	if err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *GormRefreshTokenRepository) FindByUser(userID uint) ([]RefreshToken, error) {
	var tokens []RefreshToken
	err := r.db.Where("user_id = ? AND revoked_at IS NULL", userID).Order("created_at DESC").Find(&tokens).Error
	return tokens, err
}

func (r *GormRefreshTokenRepository) RevokeByID(id uint, replacedByTokenID *uint, reason string) error {
	now := time.Now()
	updates := map[string]any{
		"revoked_at":           now,
		"revoke_reason":        reason,
		"replaced_by_token_id": replacedByTokenID,
		"updated_at":           now,
	}
	result := r.db.Model(&RefreshToken{}).Where("id = ? AND revoked_at IS NULL", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GormRefreshTokenRepository) RevokeFamily(userID uint, familyID string, reason string) error {
	now := time.Now()
	updates := map[string]any{
		"revoked_at":    now,
		"revoke_reason": reason,
		"updated_at":    now,
	}
	return r.db.Model(&RefreshToken{}).
		Where("user_id = ? AND family_id = ? AND revoked_at IS NULL", userID, familyID).
		Updates(updates).Error
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

func (r *GormRefreshTokenRepository) DeleteByUserAndDevice(userID uint, deviceID string) error {
	return r.db.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) DeleteByUserExceptDevice(userID uint, deviceID string) error {
	query := r.db.Where("user_id = ?", userID)
	if deviceID != "" {
		query = query.Where("device_id <> ?", deviceID)
	}
	return query.Delete(&RefreshToken{}).Error
}

func (r *GormRefreshTokenRepository) WithTx(tx *gorm.DB) RefreshTokenRepository {
	return &GormRefreshTokenRepository{db: tx}
}
