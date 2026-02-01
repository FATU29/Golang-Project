package handler

import (
	"Authentication_Service/internal/constant"
	_interface "Authentication_Service/internal/controller/interface"

	"github.com/gin-gonic/gin"
)

type AuthenticationHandler struct {
	AuthenticationController _interface.IAuthenticationController
}

func NewAuthenticationHandler(authenticationController _interface.IAuthenticationController) *AuthenticationHandler {
	return &AuthenticationHandler{
		AuthenticationController: authenticationController,
	}
}

func (handler *AuthenticationHandler) RouterList(r *gin.RouterGroup) {
	auth := r.Group(constant.Auth)
	auth.POST(constant.Register, handler.AuthenticationController.RegisterController)
}
