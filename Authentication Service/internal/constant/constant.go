package constant

const (
	MySql                         = "mysql"
	ACTIVE                        = 1
	NO_ACTIVE                     = 0
	ACCESS_TOKEN_EXPIRE           = 60
	REFRESH_TOKEN_EXPIRE          = 1440
	AccessTokenBlacklistKeyPrefix = "access_blacklist:"
	SSOStateKeyPrefix             = "sso_state:"

	// GoogleSSOPasswordPlaceholder is stored for users who sign in only via Google SSO (no password).
	GoogleSSOPasswordPlaceholder = "__GOOGLE_SSO__"
)
