package apperror

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

var (
	// --- 400 - BAD REQUEST ERRORS ---

	// Request & Validation Errors
	ErrBadRequest   = NewAppError("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrInvalidInput = NewAppError("INVALID_INPUT", "Invalid input provided", http.StatusBadRequest)
	ErrInvalidParam = NewAppError("INVALID_PARAM", "Invalid ID parameter", http.StatusBadRequest)
	ErrInvalidSlug  = NewAppError("INVALID_SLUG", "Invalid station slug provided", http.StatusBadRequest)
	ErrInvalidEmail = NewAppError("INVALID_EMAIL", "Invalid email format", http.StatusBadRequest)

	// Game Domain - Business Logic Errors
	ErrInsufficientCoins      = NewAppError("INSUFFICIENT_COINS", "Insufficient coins to complete this action", http.StatusBadRequest)
	ErrStationAlreadyUnlocked = NewAppError("STATION_ALREADY_UNLOCKED", "Station is already unlocked", http.StatusBadRequest)
	ErrAlreadyClaimed         = NewAppError("ALREADY_CLAIMED", "Daily reward already claimed today", http.StatusBadRequest)
	ErrUnknownRewardType      = NewAppError("UNKNOWN_REWARD_TYPE", "Unknown reward type", http.StatusBadRequest)
	ErrUserNotStartedGame     = NewAppError("USER_NOT_STARTED_GAME", "User has not started the game", http.StatusBadRequest)

	// --- 401 - UNAUTHORIZED ERRORS ---

	ErrUnauthorized      = NewAppError("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized)
	ErrInvalidToken      = NewAppError("INVALID_TOKEN", "Invalid or expired token", http.StatusUnauthorized)
	ErrTokenExpired      = NewAppError("TOKEN_EXPIRED", "Token has expired", http.StatusUnauthorized)
	ErrTokenRevoked      = NewAppError("TOKEN_REVOKED", "Token has been revoked", http.StatusUnauthorized)
	ErrMissingAuthHeader = NewAppError("MISSING_AUTH_HEADER", "Missing authorization header", http.StatusUnauthorized)

	// --- 403 - FORBIDDEN ERRORS ---

	ErrAccessDenied       = NewAppError("ACCESS_DENIED", "You don't have permission to access this resource", http.StatusForbidden)
	ErrStationNotUnlocked = NewAppError("STATION_NOT_UNLOCKED", "Station must be unlocked before upgrading", http.StatusForbidden)
	ErrStageNotUnlocked   = NewAppError("STAGE_NOT_UNLOCKED", "Stage is locked", http.StatusForbidden)

	// --- 404 - NOT FOUND ERRORS ---

	ErrRecordNotFound   = NewAppError("RECORD_NOT_FOUND", "Record not found", http.StatusNotFound)
	ErrUserNotFound     = NewAppError("USER_NOT_FOUND", "User not found", http.StatusNotFound)
	ErrStationNotFound  = NewAppError("STATION_NOT_FOUND", "Kitchen station not found", http.StatusNotFound)
	ErrFoodItemNotFound = NewAppError("FOOD_ITEM_NOT_FOUND", "Food item not found", http.StatusNotFound)
	ErrStageNotFound    = NewAppError("STAGE_NOT_FOUND", "Game stage not found", http.StatusNotFound)

	// --- 409 - CONFLICT ERRORS ---

	ErrConflict              = NewAppError("CONFLICT", "Resource already exists in the system", http.StatusConflict)
	ErrAlreadyExists         = NewAppError("ALREADY_EXISTS", "Resource already exists", http.StatusConflict)
	ErrInvalidState          = NewAppError("INVALID_STATE", "Operation cannot be performed in current state", http.StatusConflict)
	ErrMaxLevelReached       = NewAppError("MAX_LEVEL_REACHED", "Station has already reached maximum level", http.StatusConflict)
	ErrStageAlreadyCompleted = NewAppError("STAGE_ALREADY_COMPLETED", "Stage already completed", http.StatusConflict)

	// --- 500 - INTERNAL SERVER ERRORS ---

	ErrInternalServer    = NewAppError("INTERNAL_SERVER_ERROR", "An unexpected error occurred", http.StatusInternalServerError)
	ErrDatabaseError     = NewAppError("DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError)
	ErrTransactionFailed = NewAppError("TRANSACTION_FAILED", "Transaction failed", http.StatusInternalServerError)
	ErrNoUpdateRecord    = NewAppError("NO_UPDATE_RECORD", "No record found to update", http.StatusInternalServerError)
	ErrFailedRetrieveID  = NewAppError("FAILED_RETRIEVE_ID", "Failed to retrieve last inserted ID", http.StatusInternalServerError)
	ErrRequiredActiveTx  = NewAppError("REQUIRED_ACTIVE_TX", "This method requires an active transaction", http.StatusInternalServerError)
)

// AppError represents a structured application error
type AppError struct {
	Code       string // Error code for client
	Message    string // Human-readable message
	StatusCode int    // HTTP status code
	Details    string // Additional details (optional)
	Err        error  // Underlying error (optional)
}

// NewAppError creates a new application error
func NewAppError(code, message string, statusCode int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// WithDetails adds additional details to the error
func (e *AppError) WithDetails(details string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Details:    details,
		Err:        e.Err,
	}
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Details:    e.Details,
		Err:        err,
	}
}

// Is checks if the error matches the target error
func (e *AppError) Is(target error) bool {
	var t *AppError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func FiberErrorHandler(handler *ErrorHandler) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		statusCode, message, code, details := handler.HandleError(err)

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"message": message,
			"data":    nil,
			"error": fiber.Map{
				"code":    code,
				"details": details,
			},
			"timestamp": time.Now().Unix(),
		})
	}
}
