package internal

import (
	"fmt"
	"strconv"
)

func Parse[T any](value string) (T, error) {
	var zero T

	switch any(zero).(type) {
	case string:
		return any(value).(T), nil
	case int:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return zero, err
		}

		return any(parsed).(T), nil
	case bool:
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return zero, err
		}

		return any(parsed).(T), nil
	default:
		return zero, fmt.Errorf("unsupported type: %T", zero)
	}
}
