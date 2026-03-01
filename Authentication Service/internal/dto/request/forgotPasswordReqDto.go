package request

type ForgotPasswordReqDto struct {
	Email string `json:"email" validate:"required,email"`
}
