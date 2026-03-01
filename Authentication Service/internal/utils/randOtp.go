package utils

import (
	"crypto/rand"
	"io"
)

func GenerateRandomOtp(length int) (string, error) {
	const digits = "0123456789"

	otp := make([]byte, length)

	if _, err := io.ReadFull(rand.Reader, otp); err != nil {
		return "", err
	}

	for i := 0; i < length; i++ {
		otp[i] = digits[otp[i]%byte(len(digits))]
	}

	return string(otp), nil
}
