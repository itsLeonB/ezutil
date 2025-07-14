package ezutil

import (
	"fmt"
	"time"
)

// GetTimeRangeClause generates a SQL WHERE clause for time range filtering.
// It handles various combinations of start and end times, including open-ended ranges.
// Returns the SQL clause string and corresponding parameter values for prepared statements.
func GetTimeRangeClause(timeCol string, start, end time.Time) (string, []any) {
	if start.IsZero() && end.IsZero() {
		return "", nil
	}

	if start.IsZero() {
		return fmt.Sprintf("%s <= ?", timeCol), []any{end}
	}

	if end.IsZero() {
		return fmt.Sprintf("%s >= ?", timeCol), []any{start}
	}

	return fmt.Sprintf("%s BETWEEN ? AND ?", timeCol), []any{start, end}
}
