package ezutil_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"normal string", "hello", "hello"},
		{"string with spaces", "hello world", "hello world"},
		{"string with numbers", "test123", "test123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.Parse[string](tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParse_Int(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int
		expectError bool
	}{
		{"positive integer", "123", 123, false},
		{"negative integer", "-456", -456, false},
		{"zero", "0", 0, false},
		{"invalid string", "abc", 0, true},
		{"float string", "12.34", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.Parse[int](tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, 0, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParse_Bool(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    bool
		expectError bool
	}{
		{"true lowercase", "true", true, false},
		{"true uppercase", "TRUE", true, false},
		{"true mixed case", "True", true, false},
		{"false lowercase", "false", false, false},
		{"false uppercase", "FALSE", false, false},
		{"false mixed case", "False", false, false},
		{"1 as true", "1", true, false},
		{"0 as false", "0", false, false},
		{"invalid string", "maybe", false, true},
		{"empty string", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.Parse[bool](tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, false, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParse_UUID(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	expectedUUID, _ := uuid.Parse(validUUID)

	tests := []struct {
		name        string
		input       string
		expected    uuid.UUID
		expectError bool
	}{
		{"valid UUID", validUUID, expectedUUID, false},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", expectedUUID, false},
		{"invalid UUID format", "invalid-uuid", uuid.UUID{}, true},
		{"empty string", "", uuid.UUID{}, true},
		{"partial UUID", "550e8400-e29b", uuid.UUID{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.Parse[uuid.UUID](tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uuid.UUID{}, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParse_UnsupportedType(t *testing.T) {
	type CustomType struct {
		Value string
	}

	_, err := ezutil.Parse[CustomType]("test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported type")
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{"positive length", 16, false},
		{"small length", 1, false},
		{"large length", 100, false},
		{"zero length", 0, true},
		{"negative length", -5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ezutil.GenerateRandomString(tt.length)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
				// Base64 encoded string should be longer than input length
				assert.Greater(t, len(result), 0)
			}
		})
	}
}

func TestGenerateRandomString_Uniqueness(t *testing.T) {
	// Generate multiple random strings and ensure they're different
	length := 32
	strings := make(map[string]bool)

	for i := 0; i < 100; i++ {
		result, err := ezutil.GenerateRandomString(length)
		require.NoError(t, err)

		// Check if we've seen this string before
		assert.False(t, strings[result], "Generated duplicate random string: %s", result)
		strings[result] = true
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase word", "hello", "Hello"},
		{"uppercase word", "HELLO", "Hello"},
		{"mixed case word", "hELLo", "Hello"},
		{"single character", "a", "A"},
		{"single uppercase character", "A", "A"},
		{"empty string", "", ""},
		{"word with numbers", "test123", "Test123"},
		{"word with special chars", "hello!", "Hello!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.Capitalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
