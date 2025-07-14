package ezutil

import "github.com/google/uuid"

func CompareUUID(a, b uuid.UUID) int {
	for i := range len(a) {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return 0
}
