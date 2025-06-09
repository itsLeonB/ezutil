package ezutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rotisserie/eris"
)

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

func BindRequest[T any](ctx *gin.Context, bindType binding.Binding) (T, error) {
	var zero T

	if err := ctx.ShouldBindWith(&zero, bindType); err != nil {
		return zero, eris.Wrapf(err, "failed to bind request with type %s", bindType.Name())
	}

	return zero, nil
}

func GetFromContext[T any](ctx *gin.Context, key string) (T, error) {
	var zero T

	val, exists := ctx.Get(key)
	if !exists {
		return zero, eris.Errorf("value with key %s not found in context", key)
	}

	asserted, ok := val.(T)
	if !ok {
		return zero, eris.Errorf("error asserting value %s as type %T", val, zero)
	}

	return asserted, nil
}
