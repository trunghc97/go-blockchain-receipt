package api

import "github.com/labstack/echo/v4"

// HTTPError represents an error response for Swagger docs
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// Convert echo.HTTPError to our HTTPError for Swagger
func convertEchoError(err *echo.HTTPError) *HTTPError {
	msg := ""
	if m, ok := err.Message.(string); ok {
		msg = m
	} else {
		msg = "Internal Server Error"
	}
	return &HTTPError{
		Code:    err.Code,
		Message: msg,
	}
}
