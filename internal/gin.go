package internal

import (
	"crypto/subtle"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil/config"
	"github.com/rotisserie/eris"
)

func ExtractToken(ctx *gin.Context, authStrategy string) (string, string, error) {
	switch authStrategy {
	case "Bearer":
		token, errMsg := extractBearerToken(ctx)
		return token, errMsg, nil
	default:
		return "", "", eris.Errorf("unsupported auth strategy: %s", authStrategy)
	}
}

func extractBearerToken(ctx *gin.Context) (string, string) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		return "", config.MsgAuthMissingToken
	}

	isValid, token := validateAndExtractBearerToken(token)
	if !isValid {
		return "", config.MsgAuthInvalidToken
	}

	return token, ""
}

func validateAndExtractBearerToken(bearerToken string) (bool, string) {
	splits := strings.Split(bearerToken, " ")

	if len(splits) != 2 {
		return false, ""
	}

	ok := subtle.ConstantTimeCompare([]byte(splits[0]), []byte("Bearer")) == 1
	if !ok {
		return false, ""
	}

	return true, splits[1]
}
