package mapper

import (
	request "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/model"
)

func UserFromUserResDto(userResDto *model.User) *request.UserResDto {
	return &request.UserResDto{
		Id:        userResDto.Id,
		Email:     userResDto.Email,
		Lastname:  userResDto.Lastname,
		Firstname: userResDto.Firstname,
		Avatar:    userResDto.Avatar,
	}
}
