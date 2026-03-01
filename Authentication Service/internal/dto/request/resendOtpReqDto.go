package request

type ResendOtpReqDto struct {
	Email string `json:"email" validate:"required,email"`
}
