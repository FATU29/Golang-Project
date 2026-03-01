package request

type ChangePasswordReqDto struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" binding:"required,secure_password"`
}
