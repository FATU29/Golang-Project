package request

type RegisterResDto struct {
	Id        string  `json:"id"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Avatar    *string `json:"avatar"`
	Email     string  `json:"email"`
}
