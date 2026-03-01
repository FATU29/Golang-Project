package request

type IntrospectReqDto struct {
	Token         string  `json:"token" binding:"required"`
	TokenTypeHint *string `json:"token_type_hint,omitempty"` // e.g. "access_token"
}
