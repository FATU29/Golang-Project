package handler

import (
	"Authentication_Service/internal/constant"
	_interface "Authentication_Service/internal/controller/interface"
	"Authentication_Service/internal/middleware"
	serviceInterface "Authentication_Service/internal/service/interface"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	AuthenticationController _interface.IAuthenticationController
	TokenGenerator           serviceInterface.ITokenGenerator
}

func NewAuthenticationHandler(authenticationController _interface.IAuthenticationController, tokenGenerator serviceInterface.ITokenGenerator) *AuthenticationHandler {
	return &AuthenticationHandler{
		AuthenticationController: authenticationController,
		TokenGenerator:           tokenGenerator,
	}
}

func (handler *AuthenticationHandler) RouterList(r *gin.RouterGroup) {
	auth := r.Group(constant.Auth)
	auth.POST(constant.Register, handler.AuthenticationController.RegisterController)
	auth.POST(constant.ValidateMail, handler.AuthenticationController.ValidateMailController)
	auth.POST(constant.Login, handler.AuthenticationController.LoginController)
	auth.POST(constant.Logout, handler.AuthenticationController.LogoutController)
	auth.POST(constant.ResendOtp, handler.AuthenticationController.ResendOtpController)
	auth.POST(constant.ForgotPassword, handler.AuthenticationController.ForgotPasswordController)
	auth.POST(constant.ResetPassword, handler.AuthenticationController.ResetPasswordController)
	auth.POST(constant.Introspect, handler.AuthenticationController.IntrospectController)
	auth.GET(constant.GoogleSSO, handler.AuthenticationController.GoogleSSORedirectController)
	auth.GET(constant.GoogleSSOCallback, handler.AuthenticationController.GoogleSSOCallbackController)

	authProtected := auth.Group("", middleware.RequireAuth(handler.TokenGenerator))
	authProtected.POST(constant.ChangePassword, handler.AuthenticationController.ChangePasswordController)
}
