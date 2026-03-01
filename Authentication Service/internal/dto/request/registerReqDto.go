package request

type RegisterReqDto struct {
	Firstname       *string `json:"firstname"`
	Lastname        *string `json:"lastname"`
	Email           string  `json:"email" binding:"required,email"`
	Password        string  `json:"password" binding:"required,secure_password"`
	ConfirmPassword string  `json:"confirmPassword" binding:"required,eqfield=Password"`
}
