package ezutil_test

import (
	"net/http"
	"testing"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	appErr := ezutil.AppError{
		Type:           "TestError",
		Message:        "Test message",
		HttpStatusCode: http.StatusBadRequest,
		Details:        "Test details",
	}

	expected := "[TestError] Test message: Test details"
	assert.Equal(t, expected, appErr.Error())
}

func TestAppError_ErrorWithNilDetails(t *testing.T) {
	appErr := ezutil.AppError{
		Type:           "TestError",
		Message:        "Test message",
		HttpStatusCode: http.StatusBadRequest,
		Details:        nil,
	}

	expected := "[TestError] Test message: %!s(<nil>)"
	assert.Equal(t, expected, appErr.Error())
}

func TestInternalServerError(t *testing.T) {
	err := ezutil.InternalServerError()

	assert.Equal(t, "InternalServerError", err.Type)
	assert.Equal(t, "Undefined error occurred", err.Message)
	assert.Equal(t, http.StatusInternalServerError, err.HttpStatusCode)
	assert.Nil(t, err.Details)
}

func TestConflictError(t *testing.T) {
	details := map[string]string{"resource": "user", "id": "123"}
	err := ezutil.ConflictError(details)

	assert.Equal(t, "ConflictError", err.Type)
	assert.Equal(t, "Conflict with existing resource", err.Message)
	assert.Equal(t, http.StatusConflict, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestNotFoundError(t *testing.T) {
	details := "User not found"
	err := ezutil.NotFoundError(details)

	assert.Equal(t, "NotFoundError", err.Type)
	assert.Equal(t, "Requested resource is not found", err.Message)
	assert.Equal(t, http.StatusNotFound, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestUnauthorizedError(t *testing.T) {
	details := "Invalid token"
	err := ezutil.UnauthorizedError(details)

	assert.Equal(t, "UnauthorizedError", err.Type)
	assert.Equal(t, "Unauthorized access", err.Message)
	assert.Equal(t, http.StatusUnauthorized, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestForbiddenError(t *testing.T) {
	details := "Insufficient permissions"
	err := ezutil.ForbiddenError(details)

	assert.Equal(t, "ForbiddenError", err.Type)
	assert.Equal(t, "Forbidden access", err.Message)
	assert.Equal(t, http.StatusForbidden, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestBadRequestError(t *testing.T) {
	details := []string{"field1 is required", "field2 is invalid"}
	err := ezutil.BadRequestError(details)

	assert.Equal(t, "BadRequestError", err.Type)
	assert.Equal(t, "Request is not valid", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestUnprocessableEntityError(t *testing.T) {
	details := map[string]interface{}{
		"errors": []string{"validation failed"},
	}
	err := ezutil.UnprocessableEntityError(details)

	assert.Equal(t, "UnprocessableEntityError", err.Type)
	assert.Equal(t, "Request cannot be processed due to semantic errors", err.Message)
	assert.Equal(t, http.StatusUnprocessableEntity, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestValidationError(t *testing.T) {
	details := map[string]string{
		"email": "invalid format",
		"age":   "must be positive",
	}
	err := ezutil.ValidationError(details)

	assert.Equal(t, "ValidationError", err.Type)
	assert.Equal(t, "Failed to validate request", err.Message)
	assert.Equal(t, http.StatusUnprocessableEntity, err.HttpStatusCode)
	assert.Equal(t, details, err.Details)
}

func TestAppErrorImplementsError(t *testing.T) {
	var err error = ezutil.InternalServerError()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "InternalServerError")
}
