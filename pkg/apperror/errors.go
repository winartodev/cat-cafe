package apperror

import "errors"

var (
	// --- DATABASE & RESOURCE ERRORS ---

	// ErrRecordNotFound occurs when a requested database record does not exist.
	ErrRecordNotFound = errors.New("record not found")

	// ErrConflict occurs when the resource already exists in the system.
	ErrConflict = errors.New("resource already exists")

	// ErrNoUpdateRecord indicates that no rows were affected during an update operation.
	ErrNoUpdateRecord = errors.New("no record found to update")

	// ErrFailedRetrieveID occurs when the system cannot get the last inserted ID.
	ErrFailedRetrieveID = errors.New("failed to retrieve last inserted ID")

	// ErrRequiredActiveTx indicates that a database transaction is required for the operation.
	ErrRequiredActiveTx = errors.New("this method requires an active transaction")

	ErrMissingAuthHeader = errors.New("missing authorization header")

	ErrInvalidToken = errors.New("invalid token")

	ErrInvalidEmail = errors.New("invalid email")

	ErrTokenExpired = errors.New("token expired")

	ErrUnauthorized = errors.New("unauthorized")
	ErrTokenRevoked = errors.New("token has been revoked")

	// --- REQUEST & VALIDATION ERRORS ---

	// ErrBadRequest indicates that the server cannot process the request due to client error.
	ErrBadRequest = errors.New("bad request")

	// ErrInvalidParam occurs when the provided parameter is not in the correct format.
	ErrInvalidParam = errors.New("invalid id param")

	// --- DOMAIN: DAILY REWARD ---

	// ErrAlreadyClaimed indicates the user has already redeemed their reward for the day.
	ErrAlreadyClaimed = errors.New("daily reward already claimed")

	// ErrUnknownRewardType indicates that the provided reward type is not defined in the system.
	ErrUnknownRewardType = errors.New("unknown reward type")

	// --- DOMAIN: GAME STAGE ---

	ErrStageNotUnlocked      = errors.New("stage is locked")
	ErrStageAlreadyCompleted = errors.New("stage already completed")
	ErrUserNotStartedGame    = errors.New("user not started game")
	ErrMaxLevelReached       = errors.New("max level reached")
	ErrInsufficientCoins     = errors.New("insufficient coins")
)
