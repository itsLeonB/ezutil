package ezutil_test

import (
	"testing"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetTimeRangeClause(t *testing.T) {
	timeCol := "created_at"
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC)
	zeroTime := time.Time{}

	tests := []struct {
		name           string
		timeCol        string
		start          time.Time
		end            time.Time
		expectedClause string
		expectedArgs   []any
	}{
		{
			name:           "both start and end provided",
			timeCol:        timeCol,
			start:          startTime,
			end:            endTime,
			expectedClause: "created_at BETWEEN ? AND ?",
			expectedArgs:   []any{startTime, endTime},
		},
		{
			name:           "only start provided",
			timeCol:        timeCol,
			start:          startTime,
			end:            zeroTime,
			expectedClause: "created_at >= ?",
			expectedArgs:   []any{startTime},
		},
		{
			name:           "only end provided",
			timeCol:        timeCol,
			start:          zeroTime,
			end:            endTime,
			expectedClause: "created_at <= ?",
			expectedArgs:   []any{endTime},
		},
		{
			name:           "both times are zero",
			timeCol:        timeCol,
			start:          zeroTime,
			end:            zeroTime,
			expectedClause: "",
			expectedArgs:   nil,
		},
		{
			name:           "different column name",
			timeCol:        "updated_at",
			start:          startTime,
			end:            endTime,
			expectedClause: "updated_at BETWEEN ? AND ?",
			expectedArgs:   []any{startTime, endTime},
		},
		{
			name:           "column with table prefix",
			timeCol:        "users.created_at",
			start:          startTime,
			end:            endTime,
			expectedClause: "users.created_at BETWEEN ? AND ?",
			expectedArgs:   []any{startTime, endTime},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clause, args := ezutil.GetTimeRangeClause(tt.timeCol, tt.start, tt.end)

			assert.Equal(t, tt.expectedClause, clause)
			assert.Equal(t, tt.expectedArgs, args)
		})
	}
}

func TestGetTimeRangeClause_EdgeCases(t *testing.T) {
	timeCol := "timestamp_col"

	t.Run("same start and end time", func(t *testing.T) {
		sameTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

		clause, args := ezutil.GetTimeRangeClause(timeCol, sameTime, sameTime)

		assert.Equal(t, "timestamp_col BETWEEN ? AND ?", clause)
		assert.Equal(t, []any{sameTime, sameTime}, args)
	})

	t.Run("start after end time", func(t *testing.T) {
		startTime := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		clause, args := ezutil.GetTimeRangeClause(timeCol, startTime, endTime)

		// Function should still generate the clause even if logically incorrect
		assert.Equal(t, "timestamp_col BETWEEN ? AND ?", clause)
		assert.Equal(t, []any{startTime, endTime}, args)
	})

	t.Run("very old date", func(t *testing.T) {
		oldTime := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

		clause, args := ezutil.GetTimeRangeClause(timeCol, oldTime, time.Time{})

		assert.Equal(t, "timestamp_col >= ?", clause)
		assert.Equal(t, []any{oldTime}, args)
	})

	t.Run("future date", func(t *testing.T) {
		futureTime := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)

		clause, args := ezutil.GetTimeRangeClause(timeCol, time.Time{}, futureTime)

		assert.Equal(t, "timestamp_col <= ?", clause)
		assert.Equal(t, []any{futureTime}, args)
	})

	t.Run("empty column name", func(t *testing.T) {
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

		clause, args := ezutil.GetTimeRangeClause("", startTime, time.Time{})

		assert.Equal(t, " >= ?", clause)
		assert.Equal(t, []any{startTime}, args)
	})
}

func TestGetTimeRangeClause_WithDifferentTimezones(t *testing.T) {
	timeCol := "event_time"

	// Create times in different timezones
	utcTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	// EST is UTC-5
	est, _ := time.LoadLocation("America/New_York")
	estTime := time.Date(2024, 1, 15, 7, 0, 0, 0, est) // Same moment as utcTime

	t.Run("different timezone same moment", func(t *testing.T) {
		clause, args := ezutil.GetTimeRangeClause(timeCol, utcTime, estTime)

		assert.Equal(t, "event_time BETWEEN ? AND ?", clause)
		assert.Len(t, args, 2)

		// The actual time values should be preserved as-is
		assert.Equal(t, utcTime, args[0])
		assert.Equal(t, estTime, args[1])
	})

	t.Run("mixed timezone range", func(t *testing.T) {
		pst, _ := time.LoadLocation("America/Los_Angeles")
		pstTime := time.Date(2024, 1, 15, 4, 0, 0, 0, pst) // Same moment as utcTime

		clause, args := ezutil.GetTimeRangeClause(timeCol, utcTime, pstTime)

		assert.Equal(t, "event_time BETWEEN ? AND ?", clause)
		assert.Equal(t, []any{utcTime, pstTime}, args)
	})
}

func TestGetTimeRangeClause_Integration(t *testing.T) {
	// Test that the function can be used in a realistic scenario
	timeCol := "orders.created_at"

	// Get start and end of January 2024
	startOfMonth, _ := ezutil.GetStartOfDay(2024, 1, 1)
	endOfMonth, _ := ezutil.GetEndOfDay(2024, 1, 31)

	clause, args := ezutil.GetTimeRangeClause(timeCol, startOfMonth, endOfMonth)

	assert.Equal(t, "orders.created_at BETWEEN ? AND ?", clause)
	assert.Len(t, args, 2)
	assert.Equal(t, startOfMonth, args[0])
	assert.Equal(t, endOfMonth, args[1])

	// Verify the times are correct
	startTime := args[0].(time.Time)
	endTime := args[1].(time.Time)

	assert.Equal(t, 2024, startTime.Year())
	assert.Equal(t, time.January, startTime.Month())
	assert.Equal(t, 1, startTime.Day())
	assert.Equal(t, 0, startTime.Hour())

	assert.Equal(t, 2024, endTime.Year())
	assert.Equal(t, time.January, endTime.Month())
	assert.Equal(t, 31, endTime.Day())
	assert.Equal(t, 23, endTime.Hour())
}
