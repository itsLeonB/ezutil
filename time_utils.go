package ezutil

import (
	"time"

	"github.com/rotisserie/eris"
)

// GetStartOfDay creates a time.Time representing the start of the specified date (00:00:00 UTC).
// It validates the date parameters and returns an error for invalid dates.
// The returned time is in UTC timezone.
func GetStartOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

// GetEndOfDay creates a time.Time representing the end of the specified date (23:59:59.999999999 UTC).
// It validates the date parameters and returns an error for invalid dates.
// The returned time is in UTC timezone with maximum precision.
func GetEndOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

// FormatTimeNullable formats a time.Time using the specified layout, handling zero values gracefully.
// Returns an empty string if the time is zero (uninitialized), otherwise returns the formatted time.
// Useful for optional time fields in JSON responses and templates.
func FormatTimeNullable(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(layout)
}
