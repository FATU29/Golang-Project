package middleware

import (
	"Authentication_Service/internal/dto/common"
	response "Authentication_Service/internal/dto/response"
	_interface "Authentication_Service/internal/service/interface"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextKeyUser = "user"

// RequireAuth validates Bearer token and sets user in context. Use after router group that needs auth.
func RequireAuth(tokenGenerator _interface.ITokenGenerator) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		var token string

		if auth != "" {
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token = parts[1]
			}
		}

		// If no Bearer token, check for accessToken cookie
		if token == "" {
			if cookieToken, err := c.Cookie("accessToken"); err == nil {
				token = cookieToken
			}
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.ApiResponse[any]{
				Code:    http.StatusUnauthorized,
				Message: "authorization required",
			})
			return
		}

		user, err := tokenGenerator.ValidateToken(token)
		if err != nil {
			code := http.StatusUnauthorized
			if err.Error() == "token has expired" {
				code = http.StatusUnauthorized
			}
			c.AbortWithStatusJSON(code, common.ApiResponse[any]{
				Code:    code,
				Message: err.Error(),
			})
			return
		}
		c.Set(ContextKeyUser, user)
		c.Next()
	}
}

// GetUser returns the authenticated user from context. Call only after RequireAuth.
func GetUser(c *gin.Context) *response.UserResDto {
	v, ok := c.Get(ContextKeyUser)
	if !ok {
		return nil
	}
	u, _ := v.(*response.UserResDto)
	return u
}
