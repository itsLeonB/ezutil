package ezutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rotisserie/eris"
)

// GetPathParam extracts and parses a path parameter from the Gin context.
// It returns the parsed value of type T, a boolean indicating if the parameter exists,
// and an error if parsing fails. If the parameter does not exist, it returns the zero value with false.
// Supports parsing to string, int, bool, and UUID types.
func GetPathParam[T any](ctx *gin.Context, key string) (T, bool, error) {
	var zero T

	paramValue, exists := ctx.Params.Get(key)
	if !exists {
		return zero, false, nil
	}

	parsedValue, err := Parse[T](paramValue)
	if err != nil {
		return zero, false, eris.Wrapf(err, "failed to parse parameter '%s'", key)
	}

	return parsedValue, true, nil
}

func GetRequiredPathParam[T any](ctx *gin.Context, key string) (T, error) {
	var zero T

	paramValue, exists := ctx.Params.Get(key)
	if !exists {
		return zero, eris.Errorf("missing path param: %s", key)
	}

	parsedValue, err := Parse[T](paramValue)
	if err != nil {
		return zero, eris.Wrapf(err, "failed to parse parameter '%s'", key)
	}

	return parsedValue, nil
}

// BindRequest binds the incoming HTTP request to a struct of type T using the specified binding type.
// It supports various Gin binding types such as JSON, XML, Query, etc.
// Returns the bound struct or an error if binding fails.
func BindRequest[T any](ctx *gin.Context, bindType binding.Binding) (T, error) {
	var zero T

	if err := ctx.ShouldBindWith(&zero, bindType); err != nil {
		return zero, eris.Wrapf(err, "failed to bind request with type %s", bindType.Name())
	}

	return zero, nil
}

// GetFromContext retrieves a value from the Gin context and type-asserts it to type T.
// Returns the typed value or an error if the key does not exist or type assertion fails.
// Useful for retrieving typed data stored in context by middleware.
func GetFromContext[T any](ctx *gin.Context, key string) (T, error) {
	var zero T

	val, exists := ctx.Get(key)
	if !exists {
		return zero, eris.Errorf("value with key %s not found in context", key)
	}

	asserted, ok := val.(T)
	if !ok {
		return zero, eris.Errorf("error asserting value %v as type %T", val, zero)
	}

	return asserted, nil
}

func GetAndParseFromContext[T any](ctx *gin.Context, key string) (T, error) {
	var zero T

	asserted, err := GetFromContext[string](ctx, key)
	if err != nil {
		return zero, err
	}

	return Parse[T](asserted)
}
