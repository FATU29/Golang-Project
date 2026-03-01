package _interface

import (
	"context"
	"time"
)

type IOtpRepository interface {
	SaveOtp(ctx context.Context, email string, otp string, timeLife time.Duration) error
	DeleteOtp(ctx context.Context, email string) error
	ValidateOtp(ctx context.Context, email string, otp string) (bool, error)
	CheckTimeLife(ctx context.Context, email string) (bool, error)
	GetOtp(ctx context.Context, email string) (string, error)
	CheckExistedAndDelete(ctx context.Context, email string) (bool, error)

	// Forgot password OTP (separate key prefix to avoid clash with registration OTP)
	SaveForgotPasswordOtp(ctx context.Context, email string, otp string, timeLife time.Duration) error
	GetForgotPasswordOtp(ctx context.Context, email string) (string, error)
	DeleteForgotPasswordOtp(ctx context.Context, email string) error
	CheckForgotPasswordOtpTimeLife(ctx context.Context, email string) (bool, error)
}
