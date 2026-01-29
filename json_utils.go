package ezutil

import (
	"encoding/json"

	"github.com/itsLeonB/ungerr"
)

func Unmarshal[T any](data []byte) (T, error) {
	var zero T
	if err := json.Unmarshal(data, &zero); err != nil {
		return zero, ungerr.Wrapf(err, "error unmarshaling data to %T", zero)
	}
	return zero, nil
}
