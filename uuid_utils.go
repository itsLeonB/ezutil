package ezutil

import "github.com/google/uuid"

// CompareUUID compares two UUID values byte by byte.
// Returns -1 if a < b, 0 if a == b, and 1 if a > b.
// Useful for sorting UUIDs or implementing custom comparison logic.
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
