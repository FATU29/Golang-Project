package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashRefreshToken returns SHA256 hex hash of the token (for storing in DB, not for comparison with plain token).
func HashRefreshToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
