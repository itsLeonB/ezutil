package ezutil

import (
	"fmt"
	"net/http"
)

// AppError represents an application-specific error with HTTP context.
// It includes Type, Message, HttpStatusCode, and optional Details.
// AppError implements the error interface by returning a formatted error string.
type AppError struct {
	Type           string `json:"type"`
	Message        string `json:"message"`
	HttpStatusCode int    `json:"-"`
	Details        any    `json:"details,omitempty"`
}

// Error returns a formatted string representation of the AppError.
// It implements the error interface by formatting Type, Message, and Details.
func (ae AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", ae.Type, ae.Message, ae.Details)
}

// InternalServerError returns an AppError with HTTP 500 status.
// Use this for unexpected server-side errors that should be logged and investigated.
func InternalServerError() AppError {
	return AppError{
		Type:           "InternalServerError",
		Message:        "Undefined error occurred",
		HttpStatusCode: http.StatusInternalServerError,
	}
}

// ConflictError returns an AppError with HTTP 409 status.
// Use this when a request conflicts with the current state of the server,
// such as attempting to create a duplicate resource.
func ConflictError(details any) AppError {
	return AppError{
		Type:           "ConflictError",
		Message:        "Conflict with existing resource",
		HttpStatusCode: http.StatusConflict,
		Details:        details,
	}
}

// NotFoundError returns an AppError with HTTP 404 status.
// Use this when a requested resource cannot be found.
func NotFoundError(details any) AppError {
	return AppError{
		Type:           "NotFoundError",
		Message:        "Requested resource is not found",
		HttpStatusCode: http.StatusNotFound,
		Details:        details,
	}
}

// UnauthorizedError returns an AppError with HTTP 401 status.
// Use this when authentication is required but missing or invalid.
func UnauthorizedError(details any) AppError {
	return AppError{
		Type:           "UnauthorizedError",
		Message:        "Unauthorized access",
		HttpStatusCode: http.StatusUnauthorized,
		Details:        details,
	}
}

// ForbiddenError returns an AppError with HTTP 403 status.
// Use this when an authenticated user lacks permission for the requested action.
func ForbiddenError(details any) AppError {
	return AppError{
		Type:           "ForbiddenError",
		Message:        "Forbidden access",
		HttpStatusCode: http.StatusForbidden,
		Details:        details,
	}
}

// BadRequestError returns an AppError with HTTP 400 status.
// Use this when the client sends a malformed or invalid request.
func BadRequestError(details any) AppError {
	return AppError{
		Type:           "BadRequestError",
		Message:        "Request is not valid",
		HttpStatusCode: http.StatusBadRequest,
		Details:        details,
	}
}

// UnprocessableEntityError returns an AppError with HTTP 422 status.
// Use this when a well-formed request cannot be processed due to semantic errors.
func UnprocessableEntityError(details any) AppError {
	return AppError{
		Type:           "UnprocessableEntityError",
		Message:        "Request cannot be processed due to semantic errors",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}

// ValidationError returns an AppError with HTTP 422 status.
// Use this for input validation failures, providing details about the invalid fields.
func ValidationError(details any) AppError {
	return AppError{
		Type:           "ValidationError",
		Message:        "Failed to validate request",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}