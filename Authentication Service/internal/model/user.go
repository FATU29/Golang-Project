package model

type User struct {
	BaseModel
	Email     string  `json:"email"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Password  string  `json:"password"`
	Avatar    *string `json:"avatar"`
}
