package _interface

import (
	"Authentication_Service/internal/model"
	"context"
	"time"
)

type IAuthenticationRepository interface {
	FindByEmail(email string) (*model.User, error)
	SaveRefreshToken(ctx context.Context, record *model.RefreshToken) error
	AddAccessTokenToBlacklist(ctx context.Context, accessToken string, ttl time.Duration) error
	SetSSOState(ctx context.Context, state string, ttl time.Duration) error
	ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error)
}
