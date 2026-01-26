package apperror

import (
	"errors"
	"net/http"
)

type ErrorHandler struct {
	// Map of error codes to custom handlers
	customHandlers map[string]func(*AppError) (int, string, string)
}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		customHandlers: make(map[string]func(*AppError) (int, string, string)),
	}
}

// RegisterCustomHandler registers a custom handler for specific error code
func (h *ErrorHandler) RegisterCustomHandler(code string, handler func(*AppError) (int, string, string)) {
	h.customHandlers[code] = handler
}

// HandleError processes an error and returns HTTP status, message, and error code
func (h *ErrorHandler) HandleError(err error) (statusCode int, message string, code string, details string) {
	if err == nil {
		return http.StatusOK, "Success", "", ""
	}

	// Check if it's an AppError
	var appErr *AppError
	if errors.As(err, &appErr) {
		// Check for custom handler
		if handler, exists := h.customHandlers[appErr.Code]; exists {
			statusCode, message, code = handler(appErr)
			return statusCode, message, code, appErr.Details
		}

		// Default handling
		return appErr.StatusCode, appErr.Message, appErr.Code, appErr.Details
	}

	// Unknown error - treat as internal server error
	return http.StatusInternalServerError, "An unexpected error occurred", "INTERNAL_SERVER_ERROR", err.Error()
}
