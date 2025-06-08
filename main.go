package ezutil

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/itsLeonB/ezutil/types"
)

func GetParam[T any](ctx *gin.Context, paramType types.ParamType, key string) (T, bool, error) {
	var zero T

	paramValue, exists := internal.GetParamByType(ctx, paramType, key)
	if !exists {
		return zero, false, nil
	}

	parsedValue, err := internal.Parse[T](paramValue)
	if err != nil {
		return zero, false, fmt.Errorf("failed to parse parameter '%s': %w", key, err)
	}

	return parsedValue, true, nil
}
