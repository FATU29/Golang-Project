package request

// IntrospectResDto is RFC 7662-style token introspect response.
type IntrospectResDto struct {
	Active    bool    `json:"active"`
	Sub       string  `json:"sub,omitempty"`       // user id when active
	Exp       int64   `json:"exp,omitempty"`       // expiry unix seconds
	Iat       int64   `json:"iat,omitempty"`       // issued at unix seconds
	Email     string  `json:"email,omitempty"`     // when active
	Firstname *string `json:"firstname,omitempty"` // when active
	Lastname  *string `json:"lastname,omitempty"`  // when active
	Avatar     *string `json:"avatar,omitempty"`      // when active
	CoverImage *string `json:"cover_image,omitempty"` // when active
}
