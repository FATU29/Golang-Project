package common

// ApiResponse wraps all API responses with code, message, and optional data.
type ApiResponse[T any] struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"Success"`
	Data    T      `json:"data,omitempty"`
}

// ApiErrorResponse is used for error responses (no data). Use for Swagger @Failure.
type ApiErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad Request"`
}
