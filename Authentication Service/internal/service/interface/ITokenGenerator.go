package _interface

import (
	response "Authentication_Service/internal/dto/response"
	"time"
)

// IntrospectResult is returned by Introspect for token introspect API.
type IntrospectResult struct {
	User *response.UserResDto
	Exp  time.Time
	Iat  time.Time
}

type ITokenGenerator interface {
	GenerateToken(user response.UserResDto, expireMinutes int64) (string, error)
	ValidateToken(tokenString string) (*response.UserResDto, error)
	GetExpiresAt(tokenString string) (time.Time, error)
	Introspect(tokenString string) (*IntrospectResult, error)
}

