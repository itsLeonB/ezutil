package ezutil

import (
	"time"

	"github.com/rotisserie/eris"
)

func GetStartOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

func GetEndOfDay(year int, month int, day int) (time.Time, error) {
	t := time.Date(year, time.Month(month), day, 23, 59, 59, 999999999, time.UTC)
	// time.Date normalizes invalid dates, so check if the date changed
	if t.Year() != year || int(t.Month()) != month || t.Day() != day {
		return time.Time{}, eris.Errorf("invalid date: %d-%02d-%02d", year, month, day)
	}

	return t, nil
}

func FormatTimeNullable(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}

	return t.Format(layout)
}
