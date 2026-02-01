package _interface

import (
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
)

type IAuthenticationService interface {
	Login(user *request.LoginReqDto) response.LoginResDto
	Register(user *request.RegisterReqDto) response.RegisterResDto
}
