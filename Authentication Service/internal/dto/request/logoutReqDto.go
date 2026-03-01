package request

type LogoutReqDto struct {
	RefreshToken string  `json:"refreshToken" validate:"required"`
	AccessToken  *string `json:"accessToken,omitempty"`
}
