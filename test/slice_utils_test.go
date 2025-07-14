package ezutil_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapSlice(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		expected := []string{"1", "2", "3", "4", "5"}
		
		result := ezutil.MapSlice(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("string to int length", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		expected := []int{5, 5, 4}
		
		result := ezutil.MapSlice(input, func(s string) int {
			return len(s)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		expected := []string{}
		
		result := ezutil.MapSlice(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("single element", func(t *testing.T) {
		input := []int{42}
		expected := []string{"42"}
		
		result := ezutil.MapSlice(input, func(i int) string {
			return strconv.Itoa(i)
		})
		
		assert.Equal(t, expected, result)
	})

	t.Run("complex transformation", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		
		input := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 35},
		}
		
		expected := []string{"Alice (30)", "Bob (25)", "Charlie (35)"}
		
		result := ezutil.MapSlice(input, func(p Person) string {
			return p.Name + " (" + strconv.Itoa(p.Age) + ")"
		})
		
		assert.Equal(t, expected, result)
	})
}

func TestMapSliceWithError(t *testing.T) {
	t.Run("successful transformation", func(t *testing.T) {
		input := []string{"1", "2", "3", "4", "5"}
		expected := []int{1, 2, 3, 4, 5}
		
		result, err := ezutil.MapSliceWithError(input, func(s string) (int, error) {
			return strconv.Atoi(s)
		})
		
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("transformation with error", func(t *testing.T) {
		input := []string{"1", "2", "invalid", "4", "5"}
		
		result, err := ezutil.MapSliceWithError(input, func(s string) (int, error) {
			return strconv.Atoi(s)
		})
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []string{}
		expected := []int{}
		
		result, err := ezutil.MapSliceWithError(input, func(s string) (int, error) {
			return strconv.Atoi(s)
		})
		
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("error on first element", func(t *testing.T) {
		input := []string{"invalid", "2", "3"}
		
		result, err := ezutil.MapSliceWithError(input, func(s string) (int, error) {
			return strconv.Atoi(s)
		})
		
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error on last element", func(t *testing.T) {
		input := []string{"1", "2", "invalid"}
		
		result, err := ezutil.MapSliceWithError(input, func(s string) (int, error) {
			return strconv.Atoi(s)
		})
		
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("custom error", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		customErr := errors.New("custom error")
		
		result, err := ezutil.MapSliceWithError(input, func(i int) (string, error) {
			if i == 3 {
				return "", customErr
			}
			return strconv.Itoa(i), nil
		})
		
		assert.Error(t, err)
		assert.Equal(t, customErr, err)
		assert.Nil(t, result)
	})

	t.Run("complex type transformation", func(t *testing.T) {
		type Input struct {
			Value string
		}
		
		type Output struct {
			Number int
		}
		
		input := []Input{
			{Value: "10"},
			{Value: "20"},
			{Value: "30"},
		}
		
		expected := []Output{
			{Number: 10},
			{Number: 20},
			{Number: 30},
		}
		
		result, err := ezutil.MapSliceWithError(input, func(in Input) (Output, error) {
			num, err := strconv.Atoi(in.Value)
			if err != nil {
				return Output{}, err
			}
			return Output{Number: num}, nil
		})
		
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
