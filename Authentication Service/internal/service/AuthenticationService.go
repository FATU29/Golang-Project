package service

import (
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
)

type AuthenticationService struct {
}

func (a AuthenticationService) Login(user *request.LoginReqDto) response.LoginResDto {
	//TODO implement me
	panic("implement me")
}

func (a AuthenticationService) Register(user *request.RegisterReqDto) response.RegisterResDto {
	//TODO implement me
	panic("implement me")
}
