package model

import "time"

type RefreshToken struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    string    `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	DeviceInfo string   `json:"device_info"`
	IpAddress string    `json:"ip_address"`
	IsRevoked bool      `json:"is_revoked"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
