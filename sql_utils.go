package ezutil

import (
	"fmt"
	"time"
)

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
