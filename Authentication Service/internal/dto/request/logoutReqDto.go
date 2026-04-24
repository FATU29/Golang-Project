package request

type LogoutReqDto struct {
	RefreshToken string  `json:"refreshToken"`
	AccessToken  *string `json:"accessToken,omitempty"`
}
