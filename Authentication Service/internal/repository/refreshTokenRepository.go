package repository

import (
	"Authentication_Service/internal/model"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	Db *gorm.DB
}

func (r *RefreshTokenRepository) Create(refreshToken *model.RefreshToken) error {
	status := r.Db.Create(refreshToken)

	if status.Error != nil {
		return status.Error
	}

	return nil
}

func (r *RefreshTokenRepository) Delete(refreshToken *model.RefreshToken) error {
	status := r.Db.Delete(refreshToken)

	if status.Error != nil {
		return status.Error
	}

	return nil
}

func (r *RefreshTokenRepository) Update(refreshToken *model.RefreshToken) error {
	status := r.Db.Save(refreshToken)

	if status.Error != nil {
		return status.Error
	}

	return nil
}

func (r *RefreshTokenRepository) GetOne(userId string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	status := r.Db.Where("user_id = ?", userId).First(&refreshToken)

	if status.Error != nil {
		return nil, status.Error
	}

	return &refreshToken, nil
}

func (r *RefreshTokenRepository) GetByTokenHash(tokenHash string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken

	status := r.Db.Where("token_hash = ? AND is_revoked = ?", tokenHash, false).First(&refreshToken)

	if status.Error != nil {
		return nil, status.Error
	}

	return &refreshToken, nil
}

// DeleteExpired removes refresh_tokens where expires_at < now. Returns rows affected.
func (r *RefreshTokenRepository) DeleteExpired() (int64, error) {
	now := time.Now()
	result := r.Db.Where("expires_at < ?", now).Delete(&model.RefreshToken{})
	return result.RowsAffected, result.Error
}
