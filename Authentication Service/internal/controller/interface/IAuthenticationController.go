package _interface

import "github.com/gin-gonic/gin"

type IAuthenticationController interface {
	RegisterController(c *gin.Context)
	LoginController(c *gin.Context)
	LogoutController(c *gin.Context)
	ValidateMailController(c *gin.Context)
	ResendOtpController(c *gin.Context)
	ForgotPasswordController(c *gin.Context)
	ResetPasswordController(c *gin.Context)
	ChangePasswordController(c *gin.Context)
	IntrospectController(c *gin.Context)
	GoogleSSORedirectController(c *gin.Context)
	GoogleSSOCallbackController(c *gin.Context)
}
