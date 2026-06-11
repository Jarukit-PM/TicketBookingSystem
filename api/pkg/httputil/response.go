package httputil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError is the standard error payload shape for JSON responses.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorBody wraps APIError for JSON responses.
type ErrorBody struct {
	Error APIError `json:"error"`
}

// JSON writes a JSON response with the given status code.
func JSON(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

// OK writes a 200 JSON response.
func OK(c *gin.Context, body any) {
	JSON(c, http.StatusOK, body)
}

// Error writes a JSON error response.
func Error(c *gin.Context, status int, code, message string) {
	JSON(c, status, ErrorBody{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}
