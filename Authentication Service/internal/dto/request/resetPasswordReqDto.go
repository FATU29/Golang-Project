package request

type ResetPasswordReqDto struct {
	Email       string `json:"email" validate:"required,email"`
	Otp         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"newPassword" binding:"required,secure_password"`
}
