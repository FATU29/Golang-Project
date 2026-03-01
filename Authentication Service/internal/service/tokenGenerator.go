package service

import (
	"Authentication_Service/internal/config"
	response "Authentication_Service/internal/dto/response"
	"time"

	_interface "Authentication_Service/internal/service/interface"
)

type TokenGenerator struct {
	secretKey string
}

func NewTokenGenerator(secretKey string) *TokenGenerator {
	return &TokenGenerator{
		secretKey: secretKey,
	}
}

func (tg *TokenGenerator) GenerateToken(user response.UserResDto, expireMinutes int64) (string, error) {
	tokenConfig := config.TokenConfig{
		User:      user,
		SecretKey: tg.secretKey,
	}
	return tokenConfig.GenerateToken(expireMinutes)
}

func (tg *TokenGenerator) ValidateToken(tokenString string) (*response.UserResDto, error) {
	tokenConfig := config.TokenConfig{
		SecretKey: tg.secretKey,
	}
	return tokenConfig.ValidateToken(tokenString)
}

func (tg *TokenGenerator) GetExpiresAt(tokenString string) (time.Time, error) {
	tokenConfig := config.TokenConfig{
		SecretKey: tg.secretKey,
	}
	return tokenConfig.GetExpiresAt(tokenString)
}

func (tg *TokenGenerator) Introspect(tokenString string) (*_interface.IntrospectResult, error) {
	tokenConfig := config.TokenConfig{SecretKey: tg.secretKey}
	res, err := tokenConfig.GetIntrospectClaims(tokenString)
	if err != nil {
		return nil, err
	}
	return &_interface.IntrospectResult{User: res.User, Exp: res.Exp, Iat: res.Iat}, nil
}

