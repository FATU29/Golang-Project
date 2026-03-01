package config

import (
	response "Authentication_Service/internal/dto/response"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenConfig struct {
	User response.UserResDto
	jwt.RegisteredClaims
	SecretKey string
}

func (tokenCfg *TokenConfig) GenerateToken(expireMinutes int64) (string, error) {
	expirationTime := time.Duration(expireMinutes) * time.Minute

	claim := TokenConfig{
		User: tokenCfg.User,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "authentication-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(tokenCfg.SecretKey))
}

func (tokenCfg *TokenConfig) ValidateToken(tokenString string) (*response.UserResDto, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenConfig{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenCfg.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(*TokenConfig); ok && token.Valid {
		return &claims.User, nil
	}

	return nil, errors.New("could not process token claims")
}

func (tokenCfg *TokenConfig) GetClaimsFromToken(tokenString string) (*response.UserResDto, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenConfig{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenCfg.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenConfig); ok && token.Valid {
		return &claims.User, nil
	}

	return nil, errors.New("invalid token payload")
}

// GetExpiresAt parses the JWT and returns ExpiresAt (without validating exp) for computing blacklist TTL.
func (tokenCfg *TokenConfig) GetExpiresAt(tokenString string) (time.Time, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenConfig{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenCfg.SecretKey), nil
	})
	if err != nil {
		return time.Time{}, err
	}
	claims, ok := token.Claims.(*TokenConfig)
	if !ok || !token.Valid {
		return time.Time{}, errors.New("invalid token payload")
	}
	if claims.ExpiresAt == nil {
		return time.Time{}, errors.New("token has no expiration")
	}
	return claims.ExpiresAt.Time, nil
}

// IntrospectResult holds parsed token claims for introspect response.
type IntrospectResult struct {
	User *response.UserResDto
	Exp  time.Time
	Iat  time.Time
}

// GetIntrospectClaims parses and validates the token, returns user and exp/iat for introspect.
func (tokenCfg *TokenConfig) GetIntrospectClaims(tokenString string) (*IntrospectResult, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenConfig{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenCfg.SecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(*TokenConfig)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token payload")
	}
	var exp, iat time.Time
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Time
	}
	if claims.IssuedAt != nil {
		iat = claims.IssuedAt.Time
	}
	return &IntrospectResult{User: &claims.User, Exp: exp, Iat: iat}, nil
}
