package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil/types"
)

func GetParamByType(ctx *gin.Context, paramType types.ParamType, key string) (string, bool) {
	switch paramType {
	case types.ParamTypeQuery:
		return ctx.GetQuery(key)
	case types.ParamTypePath:
		return ctx.Params.Get(key)
	default:
		return "", false
	}
}
