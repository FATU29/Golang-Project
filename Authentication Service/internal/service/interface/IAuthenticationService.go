package _interface

import (
	"Authentication_Service/internal/dto/request"
	response "Authentication_Service/internal/dto/response"
	"Authentication_Service/internal/model"
	"Authentication_Service/pkg/email"
	"context"
)

type IAuthenticationService interface {
	Login(ctx context.Context, user *request.LoginReqDto) (*response.LoginResDto, *model.BusinessError)
	Register(ctx context.Context, user *request.RegisterReqDto, email email.IEmailStrategy) (*response.RegisterResDto, *model.BusinessError)
	ValidateMail(ctx context.Context, validateMailReq *request.ValidateMailReqDto) (*response.ValidateMailResDto, *model.BusinessError)
	Logout(ctx context.Context, user *request.LogoutReqDto) (*response.LogoutResDto, *model.BusinessError)
	ResendOtp(ctx context.Context, resendOtpReq *request.ResendOtpReqDto, email email.IEmailStrategy) (*response.ResendOtpResDto, *model.BusinessError)
	ForgotPassword(ctx context.Context, req *request.ForgotPasswordReqDto, emailStrategy email.IEmailStrategy) (*response.ForgotPasswordResDto, *model.BusinessError)
	ResetPassword(ctx context.Context, req *request.ResetPasswordReqDto) (*response.ResetPasswordResDto, *model.BusinessError)
	ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordReqDto) (*response.ChangePasswordResDto, *model.BusinessError)
	Introspect(ctx context.Context, req *request.IntrospectReqDto) (*response.IntrospectResDto, *model.BusinessError)
	GoogleSSO(ctx context.Context, code string) (*response.LoginResDto, *model.BusinessError)
	GetGoogleAuthURL(state string) string
	SetSSOState(ctx context.Context, state string) error
	ValidateAndConsumeSSOState(ctx context.Context, state string) (bool, error)
}
