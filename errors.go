package ezutil

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Type           string `json:"type"`
	Message        string `json:"message"`
	HttpStatusCode int    `json:"-"`
	Details        any    `json:"details,omitempty"`
}

func (ae AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", ae.Type, ae.Message, ae.Details)
}

func InternalServerError() AppError {
	return AppError{
		Type:           "InternalServerError",
		Message:        "Undefined error occurred",
		HttpStatusCode: http.StatusInternalServerError,
	}
}

func ConflictError(details any) AppError {
	return AppError{
		Type:           "ConflictError",
		Message:        "Conflict with existing resource",
		HttpStatusCode: http.StatusConflict,
		Details:        details,
	}
}

func NotFoundError(details any) AppError {
	return AppError{
		Type:           "NotFoundError",
		Message:        "Requested resource is not found",
		HttpStatusCode: http.StatusNotFound,
		Details:        details,
	}
}

func UnauthorizedError(details any) AppError {
	return AppError{
		Type:           "UnauthorizedError",
		Message:        "Unauthorized access",
		HttpStatusCode: http.StatusUnauthorized,
		Details:        details,
	}
}

func ForbiddenError(details any) AppError {
	return AppError{
		Type:           "ForbiddenError",
		Message:        "Forbidden access",
		HttpStatusCode: http.StatusForbidden,
		Details:        details,
	}
}

func BadRequestError(details any) AppError {
	return AppError{
		Type:           "BadRequestError",
		Message:        "Request is not valid",
		HttpStatusCode: http.StatusBadRequest,
		Details:        details,
	}
}

func UnprocessableEntityError(details any) AppError {
	return AppError{
		Type:           "UnprocessableEntityError",
		Message:        "Request cannot be processed due to semantic errors",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}

func ValidationError(details any) AppError {
	return AppError{
		Type:           "ValidationError",
		Message:        "Failed to validate request",
		HttpStatusCode: http.StatusUnprocessableEntity,
		Details:        details,
	}
}
