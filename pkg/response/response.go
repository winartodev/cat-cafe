package response

import (
	"github.com/gofiber/fiber/v2"
)

// response represents a standard API response format for successful requests.
type response struct {
	// Success indicates if the request was successful.
	Success bool `json:"success"`

	// Message is an optional string describing the result of the operation.
	Message string `json:"message,omitempty"`

	// Data contains the main payload of the response.
	Data interface{} `json:"data,omitempty"`

	// Meta contains additional metadata, such as pagination info.
	Meta interface{} `json:"meta,omitempty"`
}

// errorResponse represents a standard API response format for failed requests.
type errorResponse struct {
	// Success will always be false for this response type.
	Success bool `json:"success"`

	// Message describes the general category of the error.
	Message string `json:"message,omitempty"`

	// Error contains the technical error details or specific field validation errors.
	Error interface{} `json:"error,omitempty"`
}

// SuccessResponse sends a JSON response with a success status code and the provided payload.
func SuccessResponse(c *fiber.Ctx, code int, message string, data interface{}, meta interface{}) error {
	return c.Status(code).JSON(response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// FailedResponse sends a JSON response when a system or logic error occurs.
func FailedResponse(c *fiber.Ctx, code int, err error) error {
	return c.Status(code).JSON(errorResponse{
		Success: false,
		Message: "Operation failed",
		Error:   err.Error(),
	})
}
