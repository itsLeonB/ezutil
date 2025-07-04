package ezutil

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/ezutil/config"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
)

// NewCorsMiddleware creates a CORS middleware for Gin with the provided configuration.
// If corsConfig is nil, default settings are used (via cors.Default()).
// The middleware validates the configuration and logs a fatal error if invalid.
// Returns a Gin HandlerFunc to handle CORS according to the specified config.
func NewCorsMiddleware(corsConfig *cors.Config) gin.HandlerFunc {
	if corsConfig == nil {
		log.Println("CORS configuration is nil, using default settings")
		return cors.Default()
	}

	if err := corsConfig.Validate(); err != nil {
		log.Fatalf("invalid CORS configuration: %s", err.Error())
	}

	return cors.New(*corsConfig)
}

// NewAuthMiddleware creates an authentication middleware for Gin.
// It extracts a token using the given strategy (e.g., "header" or "cookie") via internal.ExtractToken,
// calls tokenCheckFunc to validate the token and retrieve user data,
// stores user data in the Gin context, and aborts the request on errors.
// Returns a Gin HandlerFunc for authentication handling.
func NewAuthMiddleware(
	authStrategy string,
	tokenCheckFunc func(ctx *gin.Context, token string) (bool, map[string]any, error),
) gin.HandlerFunc {
	if tokenCheckFunc == nil {
		log.Fatalf("tokenCheckFunc cannot be nil")
	}

	return func(ctx *gin.Context) {
		token, errMsg, err := internal.ExtractToken(ctx, authStrategy)
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, "error extracting token"))
			ctx.Abort()
			return
		}
		if errMsg != "" {
			_ = ctx.Error(UnauthorizedError(errMsg))
			ctx.Abort()
			return
		}

		exists, data, err := tokenCheckFunc(ctx, token)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}
		if !exists {
			_ = ctx.Error(UnauthorizedError(config.MsgAuthUserNotFound))
			ctx.Abort()
			return
		}

		for key, val := range data {
			ctx.Set(key, val)
		}

		ctx.Next()
	}
}

// NewPermissionMiddleware creates a permission-checking middleware for Gin.
// It retrieves the user role from context using the provided roleContextKey,
// checks if the role exists in permissionMap and includes the requiredPermission,
// and aborts the request with a ForbiddenError if permission is missing.
// Returns a Gin HandlerFunc for permission enforcement.
func NewPermissionMiddleware(
	roleContextKey string,
	requiredPermission string,
	permissionMap map[string][]string,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role := ctx.GetString(roleContextKey)
		if role == "" {
			_ = ctx.Error(eris.Errorf("role not found in context or invalid type"))
			ctx.Abort()
			return
		}

		permissions, ok := permissionMap[role]
		if !ok {
			_ = ctx.Error(eris.Errorf("unknown role: %s", role))
			ctx.Abort()
			return
		}

		if !slices.Contains(permissions, requiredPermission) {
			_ = ctx.Error(ForbiddenError(config.MsgNoPermission))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// NewErrorMiddleware creates an error handling middleware for Gin.
// It should be registered last and captures errors from previous handlers,
// converts them into AppError or validation errors, and sends a structured JSON response
// with the appropriate HTTP status code. Returns a Gin HandlerFunc.
func NewErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if err := ctx.Errors.Last(); err != nil {
			if originalErr, ok := err.Err.(AppError); ok {
				ctx.AbortWithStatusJSON(originalErr.HttpStatusCode, NewErrorResponse(originalErr))
				return
			}

			statusCode, appError := constructAppError(err)
			ctx.AbortWithStatusJSON(statusCode, NewErrorResponse(appError))
		}
	}
}

func constructAppError(err *gin.Error) (int, error) {
	originalErr := eris.Unwrap(err.Err)
	switch originalErr := originalErr.(type) {
	case validator.ValidationErrors:
		var errors []string
		for _, e := range originalErr {
			errors = append(errors, e.Error())
		}

		return http.StatusUnprocessableEntity, ValidationError(errors)
	case *json.SyntaxError:
		return http.StatusBadRequest, BadRequestError(config.MsgInvalidJson)
	default:
		// EOF error from json package is unexported
		if originalErr == io.EOF || originalErr.Error() == "EOF" {
			return http.StatusBadRequest, BadRequestError(config.MsgMissingBody)
		}

		log.Printf("unhandled error of type: %T\n", originalErr)
		log.Println(eris.ToString(err.Err, true))
		return http.StatusInternalServerError, InternalServerError()
	}
}
