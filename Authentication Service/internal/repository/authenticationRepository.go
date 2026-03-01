package repository

import (
	"Authentication_Service/internal/constant"
	"Authentication_Service/internal/model"
	"Authentication_Service/internal/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthenticationRepository struct {
	Db *gorm.DB
	Rd *redis.Client
}

func (a *AuthenticationRepository) SaveRefreshToken(ctx context.Context, record *model.RefreshToken) error {
	if record.Id == "" {
		record.Id = uuid.New().String()
	}
	return a.Db.WithContext(ctx).Create(record).Error
}

func (a *AuthenticationRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	result := a.Db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (a *AuthenticationRepository) AddAccessTokenToBlacklist(ctx context.Context, accessToken string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = time.Second
	}

	key := constant.AccessTokenBlacklistKeyPrefix + utils.HashRefreshToken(accessToken)
	return a.Rd.Set(ctx, key, "1", ttl).Err()
}

func (a *AuthenticationRepository) SetSSOState(ctx context.Context, state string, ttl time.Duration) error {
	if state == "" || ttl <= 0 {
		return errors.New("invalid state or ttl")
	}
	key := constant.SSOStateKeyPrefix + state
	return a.Rd.Set(ctx, key, "1", ttl).Err()
}

func (a *AuthenticationRepository) ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error) {
	if state == "" {
		return false, nil
	}
	key := constant.SSOStateKeyPrefix + state
	ok, err := a.Rd.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	_ = a.Rd.Del(ctx, key).Err()
	return ok == "1", nil
}
