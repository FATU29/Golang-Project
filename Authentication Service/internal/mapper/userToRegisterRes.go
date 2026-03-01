package mapper

import (
	response "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/model"
)

func FromUserModelToRegisterRes(user *model.User) *response.RegisterResDto {
	return &response.RegisterResDto{
		Id:        user.Id,
		Email:     user.Email,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
		Avatar:    user.Avatar,
	}
}
