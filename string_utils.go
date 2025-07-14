package ezutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

// Parse converts a string value to the specified type T.
// Supported types include string, int, bool, and uuid.UUID.
// Returns an error if parsing fails or the type is unsupported.
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

// GenerateRandomString creates a cryptographically secure random string of the specified length.
// The string is base64-encoded using URL-safe encoding without padding.
// Returns an error if length is non-positive or random generation fails.
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

// Capitalize converts the first character of a word to uppercase and the rest to lowercase.
// Returns an empty string if the input is empty.
// Useful for formatting names and titles consistently.
func Capitalize(word string) string {
	if len(word) == 0 {
		return ""
	}
	return string(unicode.ToUpper(rune(word[0]))) + strings.ToLower(word[1:])
}
