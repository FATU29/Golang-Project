package request

type LoginReqDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
