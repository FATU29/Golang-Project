package request

type UserResDto struct {
	Id        string  `json:"id"`
	Firstname *string `json:"firstname"`
	Lastname  *string `json:"lastname"`
	Avatar     *string `json:"avatar"`
	CoverImage *string `json:"cover_image"`
	Email      string  `json:"email"`
}
