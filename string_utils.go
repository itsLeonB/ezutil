package ezutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func Parse[T any](value string) (T, error) {
	var parsed any
	var err error
	var zero T

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		parsed, err = strconv.Atoi(value)
	case bool:
		parsed, err = strconv.ParseBool(value)
	case uuid.UUID:
		parsed, err = uuid.Parse(value)
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}

	if err != nil {
		return zero, eris.Wrapf(err, "failed to parse value '%s' as %T", value, zero)
	}

	return parsed.(T), nil
}

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", eris.New("length must be greater than 0")
	}

	randomBytes := make([]byte, length)

	_, err := io.ReadFull(rand.Reader, randomBytes)
	if err != nil {
		return "", eris.Wrap(err, "failed to generate random string")
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
