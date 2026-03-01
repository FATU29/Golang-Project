package model

type BusinessError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func NewBusinessError(statusCode int, message string) *BusinessError {
	return &BusinessError{
		StatusCode: statusCode,
		Message:    message,
	}
}
