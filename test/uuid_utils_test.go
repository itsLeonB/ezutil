package ezutil_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
)

func TestCompareUUID(t *testing.T) {
	// Create test UUIDs
	uuid1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	uuid2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	uuid3 := uuid.MustParse("00000000-0000-0000-0000-000000000001") // Same as uuid1
	uuid4 := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	zeroUUID := uuid.UUID{} // All zeros

	tests := []struct {
		name     string
		a        uuid.UUID
		b        uuid.UUID
		expected int
	}{
		{
			name:     "a < b",
			a:        uuid1,
			b:        uuid2,
			expected: -1,
		},
		{
			name:     "a > b",
			a:        uuid2,
			b:        uuid1,
			expected: 1,
		},
		{
			name:     "a == b",
			a:        uuid1,
			b:        uuid3,
			expected: 0,
		},
		{
			name:     "zero UUID vs non-zero",
			a:        zeroUUID,
			b:        uuid1,
			expected: -1,
		},
		{
			name:     "non-zero vs zero UUID",
			a:        uuid1,
			b:        zeroUUID,
			expected: 1,
		},
		{
			name:     "both zero UUIDs",
			a:        zeroUUID,
			b:        uuid.UUID{},
			expected: 0,
		},
		{
			name:     "max UUID vs normal",
			a:        uuid4,
			b:        uuid1,
			expected: 1,
		},
		{
			name:     "normal vs max UUID",
			a:        uuid1,
			b:        uuid4,
			expected: -1,
		},
		{
			name:     "identical max UUIDs",
			a:        uuid4,
			b:        uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.CompareUUID(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompareUUID_Properties(t *testing.T) {
	// Test reflexivity: CompareUUID(a, a) == 0
	t.Run("reflexivity", func(t *testing.T) {
		testUUIDs := []uuid.UUID{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			uuid.MustParse("00000000-0000-0000-0000-000000000000"),
			uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
			uuid.New(),
		}

		for _, u := range testUUIDs {
			result := ezutil.CompareUUID(u, u)
			assert.Equal(t, 0, result, "CompareUUID should return 0 for identical UUIDs")
		}
	})

	// Test antisymmetry: if CompareUUID(a, b) == x, then CompareUUID(b, a) == -x
	t.Run("antisymmetry", func(t *testing.T) {
		uuid1 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		uuid2 := uuid.MustParse("550e8400-e29b-41d4-a716-446655440001")

		result1 := ezutil.CompareUUID(uuid1, uuid2)
		result2 := ezutil.CompareUUID(uuid2, uuid1)

		assert.Equal(t, -result1, result2, "CompareUUID should be antisymmetric")
	})

	// Test transitivity: if CompareUUID(a, b) <= 0 and CompareUUID(b, c) <= 0, then CompareUUID(a, c) <= 0
	t.Run("transitivity", func(t *testing.T) {
		uuid1 := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		uuid2 := uuid.MustParse("00000000-0000-0000-0000-000000000002")
		uuid3 := uuid.MustParse("00000000-0000-0000-0000-000000000003")

		result12 := ezutil.CompareUUID(uuid1, uuid2)
		result23 := ezutil.CompareUUID(uuid2, uuid3)
		result13 := ezutil.CompareUUID(uuid1, uuid3)

		if result12 <= 0 && result23 <= 0 {
			assert.LessOrEqual(t, result13, 0, "CompareUUID should be transitive")
		}
	})
}

func TestCompareUUID_ByteByByteComparison(t *testing.T) {
	// Test that comparison works byte by byte
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{
			name:     "first byte different",
			a:        "00000000-0000-0000-0000-000000000000",
			b:        "01000000-0000-0000-0000-000000000000",
			expected: -1,
		},
		{
			name:     "last byte different",
			a:        "00000000-0000-0000-0000-000000000000",
			b:        "00000000-0000-0000-0000-000000000001",
			expected: -1,
		},
		{
			name:     "middle byte different",
			a:        "00000000-0000-0000-0000-000000000000",
			b:        "00000000-0000-0100-0000-000000000000",
			expected: -1,
		},
		{
			name:     "multiple bytes different - first difference wins",
			a:        "00000000-0000-0000-0000-000000000001",
			b:        "01000000-0000-0000-0000-000000000000",
			expected: -1, // First byte difference (00 < 01) determines result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidA := uuid.MustParse(tt.a)
			uuidB := uuid.MustParse(tt.b)

			result := ezutil.CompareUUID(uuidA, uuidB)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompareUUID_CanBeUsedForSorting(t *testing.T) {
	// Create a slice of UUIDs
	uuids := []uuid.UUID{
		uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
	}

	// Sort using CompareUUID
	for i := 0; i < len(uuids)-1; i++ {
		for j := i + 1; j < len(uuids); j++ {
			if ezutil.CompareUUID(uuids[i], uuids[j]) > 0 {
				uuids[i], uuids[j] = uuids[j], uuids[i]
			}
		}
	}

	// Verify sorted order
	expected := []uuid.UUID{
		uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
		uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
	}

	assert.Equal(t, expected, uuids)

	// Verify that each adjacent pair is in correct order
	for i := 0; i < len(uuids)-1; i++ {
		result := ezutil.CompareUUID(uuids[i], uuids[i+1])
		assert.LessOrEqual(t, result, 0, "UUIDs should be in sorted order")
	}
}
