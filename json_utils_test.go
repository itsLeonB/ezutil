package ezutil_test

import (
	"testing"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestUnmarshal(t *testing.T) {
	t.Run("valid JSON to struct", func(t *testing.T) {
		data := []byte(`{"name":"John","age":30}`)

		result, err := ezutil.Unmarshal[TestStruct](data)

		require.NoError(t, err)
		assert.Equal(t, "John", result.Name)
		assert.Equal(t, 30, result.Age)
	})

	t.Run("valid JSON to map", func(t *testing.T) {
		data := []byte(`{"key":"value","number":42}`)

		result, err := ezutil.Unmarshal[map[string]interface{}](data)

		require.NoError(t, err)
		assert.Equal(t, "value", result["key"])
		assert.Equal(t, float64(42), result["number"])
	})

	t.Run("valid JSON to slice", func(t *testing.T) {
		data := []byte(`[1,2,3,4,5]`)

		result, err := ezutil.Unmarshal[[]int](data)

		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("empty JSON object", func(t *testing.T) {
		data := []byte(`{}`)

		result, err := ezutil.Unmarshal[TestStruct](data)

		require.NoError(t, err)
		assert.Equal(t, "", result.Name)
		assert.Equal(t, 0, result.Age)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(`{"name":"John","age":}`)

		result, err := ezutil.Unmarshal[TestStruct](data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling data")
		assert.Equal(t, TestStruct{}, result)
	})

	t.Run("malformed JSON", func(t *testing.T) {
		data := []byte(`not json at all`)

		result, err := ezutil.Unmarshal[map[string]interface{}](data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling data")
		assert.Nil(t, result)
	})

	t.Run("type mismatch", func(t *testing.T) {
		data := []byte(`"string value"`)

		result, err := ezutil.Unmarshal[int](data)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error unmarshaling data")
		assert.Equal(t, 0, result)
	})
}
