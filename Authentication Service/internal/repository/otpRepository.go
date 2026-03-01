package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type OtpRepository struct {
	Rd *redis.Client
}

func (o *OtpRepository) CheckExistedAndDelete(ctx context.Context, email string) (bool, error) {
	res, err := o.Rd.Exists(ctx, email).Result()

	if err != nil {
		return false, err
	}

	if res == 1 {
		errDel := o.DeleteOtp(ctx, email)
		if errDel != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (o *OtpRepository) SaveOtp(ctx context.Context, email string, otp string, timeLife time.Duration) error {
	_, err := o.CheckExistedAndDelete(ctx, email)

	if err != nil {
		return err
	}

	return o.Rd.Set(ctx, email, otp, timeLife).Err()
}

func (o *OtpRepository) DeleteOtp(ctx context.Context, email string) error {
	return o.Rd.Del(ctx, email).Err()
}

func (o *OtpRepository) ValidateOtp(ctx context.Context, email string, otp string) (bool, error) {
	val, err := o.Rd.Get(ctx, email).Result()

	if err != nil {
		return false, err
	}

	if val != otp {
		return false, nil
	}

	ok, err := o.CheckTimeLife(ctx, email)
	return ok, err
}

func (o *OtpRepository) CheckTimeLife(ctx context.Context, email string) (bool, error) {
	ttl, err := o.Rd.TTL(ctx, email).Result()
	if err != nil {
		return false, err
	}

	if ttl == -2 {
		return false, nil
	}

	return ttl > 0 || ttl == -1, nil
}

func (o *OtpRepository) GetOtp(ctx context.Context, email string) (string, error) {
	return o.Rd.Get(ctx, email).Result()
}

const forgotPasswordKeyPrefix = "forgot_password:"

func forgotPasswordKey(email string) string {
	return forgotPasswordKeyPrefix + email
}

func (o *OtpRepository) SaveForgotPasswordOtp(ctx context.Context, email string, otp string, timeLife time.Duration) error {
	return o.Rd.Set(ctx, forgotPasswordKey(email), otp, timeLife).Err()
}

func (o *OtpRepository) GetForgotPasswordOtp(ctx context.Context, email string) (string, error) {
	return o.Rd.Get(ctx, forgotPasswordKey(email)).Result()
}

func (o *OtpRepository) DeleteForgotPasswordOtp(ctx context.Context, email string) error {
	return o.Rd.Del(ctx, forgotPasswordKey(email)).Err()
}

func (o *OtpRepository) CheckForgotPasswordOtpTimeLife(ctx context.Context, email string) (bool, error) {
	ttl, err := o.Rd.TTL(ctx, forgotPasswordKey(email)).Result()
	if err != nil {
		return false, err
	}
	if ttl == -2 {
		return false, nil
	}
	return ttl > 0 || ttl == -1, nil
}
