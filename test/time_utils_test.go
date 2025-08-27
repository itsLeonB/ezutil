package ezutil_test

import (
	"testing"
	"time"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStartOfDay(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		month       int
		day         int
		expectError bool
		expected    time.Time
	}{
		{
			name:        "valid date",
			year:        2024,
			month:       1,
			day:         15,
			expectError: false,
			expected:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "leap year date",
			year:        2024,
			month:       2,
			day:         29,
			expectError: false,
			expected:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "new year",
			year:        2024,
			month:       1,
			day:         1,
			expectError: false,
			expected:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "end of year",
			year:        2024,
			month:       12,
			day:         31,
			expectError: false,
			expected:    time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:        "invalid date - February 30",
			year:        2024,
			month:       2,
			day:         30,
			expectError: true,
		},
		{
			name:        "invalid date - month 13",
			year:        2024,
			month:       13,
			day:         1,
			expectError: true,
		},
		{
			name:        "invalid date - day 0",
			year:        2024,
			month:       1,
			day:         0,
			expectError: true,
		},
		{
			name:        "invalid date - negative year",
			year:        -1,
			month:       1,
			day:         1,
			expectError: false, // time.Date handles negative years
			expected:    time.Date(-1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.GetStartOfDay(tt.year, tt.month, tt.day)

			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, result.IsZero())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, time.UTC, result.Location())
			}
		})
	}
}

func TestGetEndOfDay(t *testing.T) {
	tests := []struct {
		name        string
		year        int
		month       int
		day         int
		expectError bool
		expected    time.Time
	}{
		{
			name:        "valid date",
			year:        2024,
			month:       1,
			day:         15,
			expectError: false,
			expected:    time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:        "leap year date",
			year:        2024,
			month:       2,
			day:         29,
			expectError: false,
			expected:    time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:        "new year",
			year:        2024,
			month:       1,
			day:         1,
			expectError: false,
			expected:    time.Date(2024, 1, 1, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:        "end of year",
			year:        2024,
			month:       12,
			day:         31,
			expectError: false,
			expected:    time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name:        "invalid date - February 30",
			year:        2024,
			month:       2,
			day:         30,
			expectError: true,
		},
		{
			name:        "invalid date - month 13",
			year:        2024,
			month:       13,
			day:         1,
			expectError: true,
		},
		{
			name:        "invalid date - day 0",
			year:        2024,
			month:       1,
			day:         0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.GetEndOfDay(tt.year, tt.month, tt.day)

			if tt.expectError {
				assert.Error(t, err)
				assert.True(t, result.IsZero())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, time.UTC, result.Location())
			}
		})
	}
}

func TestFormatTimeNullable(t *testing.T) {
	layout := "2006-01-02 15:04:05"

	tests := []struct {
		name     string
		time     time.Time
		layout   string
		expected string
	}{
		{
			name:     "valid time",
			time:     time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			layout:   layout,
			expected: "2024-01-15 14:30:45",
		},
		{
			name:     "zero time",
			time:     time.Time{},
			layout:   layout,
			expected: "",
		},
		{
			name:     "different layout",
			time:     time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			layout:   "2006/01/02",
			expected: "2024/01/15",
		},
		{
			name:     "time only layout",
			time:     time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			layout:   "15:04:05",
			expected: "14:30:45",
		},
		{
			name:     "RFC3339 layout",
			time:     time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			layout:   time.RFC3339,
			expected: "2024-01-15T14:30:45Z",
		},
		{
			name:     "custom format",
			time:     time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC),
			layout:   "Monday, January 2, 2006",
			expected: "Monday, January 15, 2024",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.FormatTimeNullable(tt.time, tt.layout)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimeUtilsIntegration(t *testing.T) {
	// Test that start and end of day work together
	year, month, day := 2024, 1, 15

	start, err := ezutil.GetStartOfDay(year, month, day)
	require.NoError(t, err)

	end, err := ezutil.GetEndOfDay(year, month, day)
	require.NoError(t, err)

	// Start should be before end
	assert.True(t, start.Before(end))

	// They should be on the same date
	assert.Equal(t, start.Year(), end.Year())
	assert.Equal(t, start.Month(), end.Month())
	assert.Equal(t, start.Day(), end.Day())

	// Start should be at midnight
	assert.Equal(t, 0, start.Hour())
	assert.Equal(t, 0, start.Minute())
	assert.Equal(t, 0, start.Second())
	assert.Equal(t, 0, start.Nanosecond())

	// End should be at end of day
	assert.Equal(t, 23, end.Hour())
	assert.Equal(t, 59, end.Minute())
	assert.Equal(t, 59, end.Second())
	assert.Equal(t, 999999999, end.Nanosecond())

	// Test formatting
	startFormatted := ezutil.FormatTimeNullable(start, "2006-01-02 15:04:05")
	endFormatted := ezutil.FormatTimeNullable(end, "2006-01-02 15:04:05")

	assert.Equal(t, "2024-01-15 00:00:00", startFormatted)
	assert.Equal(t, "2024-01-15 23:59:59", endFormatted)
}
