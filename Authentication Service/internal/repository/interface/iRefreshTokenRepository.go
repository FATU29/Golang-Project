package _interface

import "Authentication_Service/internal/model"

type IRefreshTokenRepository interface {
	Create(refreshToken *model.RefreshToken) error
	Delete(refreshToken *model.RefreshToken) error
	Update(refreshToken *model.RefreshToken) error
	GetOne(userId string) (*model.RefreshToken, error)
	GetByTokenHash(tokenHash string) (*model.RefreshToken, error)
	DeleteExpired() (int64, error)
}
