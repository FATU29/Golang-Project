package mapper

import (
	"Authentication_Service/internal/dto/request"
	res "Authentication_Service/internal/dto/response"
)

func RegisterReqToRegisterRes(req *request.RegisterReqDto) *res.RegisterResDto {
	return &res.RegisterResDto{
		Id:        "", // set by caller from created user when needed
		Email:     req.Email,
		Lastname:  req.Lastname,
		Firstname: req.Firstname,
	}
}
