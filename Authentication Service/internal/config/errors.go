package config

import (
	"Authentication_Service/internal/model"
	"net/http"
)

var (
	InternalServerError = &model.BusinessError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}
	GenOtpFailed = &model.BusinessError{
		StatusCode: http.StatusBadRequest,
		Message:    "generate otp failed",
	}
	SaveUserRedisError = &model.BusinessError{
		StatusCode: http.StatusInternalServerError,
		Message:    "save user redis error",
	}
	UserNotFound = &model.BusinessError{
		StatusCode: http.StatusNotFound,
		Message:    "user not found",
	}
	InvalidCredentials = &model.BusinessError{
		StatusCode: http.StatusUnauthorized,
		Message:    "invalid credentials",
	}
	UserAlreadyExists = &model.BusinessError{
		StatusCode: http.StatusConflict,
		Message:    "user already exists and is active",
	}
	SendMailError = &model.BusinessError{
		StatusCode: http.StatusInternalServerError,
		Message:    "send mail failed",
	}
	PayloadMailError = &model.BusinessError{
		StatusCode: http.StatusBadRequest,
		Message:    "failed to marshal JSON payload",
	}
	NetworkMailError = &model.BusinessError{
		StatusCode: http.StatusInternalServerError,
		Message:    "network error while sending mail",
	}
	GetOtpFailedError = &model.BusinessError{
		StatusCode: http.StatusBadRequest,
		Message:    "get otp failed",
	}
	InvalidOtpError = &model.BusinessError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid otp",
	}
	UserNotActivated = &model.BusinessError{
		StatusCode: http.StatusForbidden,
		Message:    "user not activated",
	}

	SessionNotFound = &model.BusinessError{
		StatusCode: http.StatusNotFound,
		Message:    "session not found",
	}
	InvalidOldPassword = &model.BusinessError{
		StatusCode: http.StatusBadRequest,
		Message:    "current password is incorrect",
	}
)
