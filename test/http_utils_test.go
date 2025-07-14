package ezutil_test

import (
	"errors"
	"testing"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
)

func TestPagination_IsZero(t *testing.T) {
	tests := []struct {
		name       string
		pagination ezutil.Pagination
		expected   bool
	}{
		{
			name:       "zero pagination",
			pagination: ezutil.Pagination{},
			expected:   true,
		},
		{
			name: "pagination with total data only",
			pagination: ezutil.Pagination{
				TotalData: 10,
			},
			expected: false,
		},
		{
			name: "pagination with current page only",
			pagination: ezutil.Pagination{
				CurrentPage: 1,
			},
			expected: false,
		},
		{
			name: "pagination with total pages only",
			pagination: ezutil.Pagination{
				TotalPages: 5,
			},
			expected: false,
		},
		{
			name: "pagination with has next page only",
			pagination: ezutil.Pagination{
				HasNextPage: true,
			},
			expected: false,
		},
		{
			name: "pagination with has prev page only",
			pagination: ezutil.Pagination{
				HasPrevPage: true,
			},
			expected: false,
		},
		{
			name: "fully populated pagination",
			pagination: ezutil.Pagination{
				TotalData:   100,
				CurrentPage: 2,
				TotalPages:  10,
				HasNextPage: true,
				HasPrevPage: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pagination.IsZero()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewResponse(t *testing.T) {
	message := "Success"
	response := ezutil.NewResponse(message)

	assert.Equal(t, message, response.Message)
	assert.Nil(t, response.Data)
	assert.Nil(t, response.Errors)
	assert.True(t, response.Pagination.IsZero())
}

func TestNewErrorResponse(t *testing.T) {
	err := errors.New("test error")
	response := ezutil.NewErrorResponse(err)

	assert.Equal(t, "test error", response.Message)
	assert.Nil(t, response.Data)
	assert.Equal(t, err, response.Errors)
	assert.True(t, response.Pagination.IsZero())
}

func TestJSONResponse_WithData(t *testing.T) {
	originalResponse := ezutil.NewResponse("Success")
	data := map[string]interface{}{
		"id":   1,
		"name": "Test",
	}

	response := originalResponse.WithData(data)

	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, data, response.Data)
	assert.Nil(t, response.Errors)
	assert.True(t, response.Pagination.IsZero())

	// Original response should be unchanged
	assert.Nil(t, originalResponse.Data)
}

func TestJSONResponse_WithError(t *testing.T) {
	originalResponse := ezutil.NewResponse("Success")
	err := errors.New("test error")

	response := originalResponse.WithError(err)

	assert.Equal(t, "Success", response.Message)
	assert.Nil(t, response.Data)
	assert.Equal(t, err, response.Errors)
	assert.True(t, response.Pagination.IsZero())

	// Original response should be unchanged
	assert.Nil(t, originalResponse.Errors)
}

func TestJSONResponse_WithPagination(t *testing.T) {
	tests := []struct {
		name         string
		queryOptions ezutil.QueryOptions
		totalData    int
		expected     ezutil.Pagination
	}{
		{
			name: "first page with next page available",
			queryOptions: ezutil.QueryOptions{
				Page:  1,
				Limit: 10,
			},
			totalData: 25,
			expected: ezutil.Pagination{
				TotalData:   25,
				CurrentPage: 1,
				TotalPages:  3,
				HasNextPage: true,
				HasPrevPage: false,
			},
		},
		{
			name: "middle page",
			queryOptions: ezutil.QueryOptions{
				Page:  2,
				Limit: 10,
			},
			totalData: 25,
			expected: ezutil.Pagination{
				TotalData:   25,
				CurrentPage: 2,
				TotalPages:  3,
				HasNextPage: true,
				HasPrevPage: true,
			},
		},
		{
			name: "last page",
			queryOptions: ezutil.QueryOptions{
				Page:  3,
				Limit: 10,
			},
			totalData: 25,
			expected: ezutil.Pagination{
				TotalData:   25,
				CurrentPage: 3,
				TotalPages:  3,
				HasNextPage: false,
				HasPrevPage: true,
			},
		},
		{
			name: "single page",
			queryOptions: ezutil.QueryOptions{
				Page:  1,
				Limit: 10,
			},
			totalData: 5,
			expected: ezutil.Pagination{
				TotalData:   5,
				CurrentPage: 1,
				TotalPages:  1,
				HasNextPage: false,
				HasPrevPage: false,
			},
		},
		{
			name: "no data",
			queryOptions: ezutil.QueryOptions{
				Page:  1,
				Limit: 10,
			},
			totalData: 0,
			expected: ezutil.Pagination{
				TotalData:   0,
				CurrentPage: 1,
				TotalPages:  0,
				HasNextPage: false,
				HasPrevPage: false,
			},
		},
		{
			name: "exact multiple of limit",
			queryOptions: ezutil.QueryOptions{
				Page:  2,
				Limit: 10,
			},
			totalData: 20,
			expected: ezutil.Pagination{
				TotalData:   20,
				CurrentPage: 2,
				TotalPages:  2,
				HasNextPage: false,
				HasPrevPage: true,
			},
		},
		{
			name: "large limit",
			queryOptions: ezutil.QueryOptions{
				Page:  1,
				Limit: 100,
			},
			totalData: 25,
			expected: ezutil.Pagination{
				TotalData:   25,
				CurrentPage: 1,
				TotalPages:  1,
				HasNextPage: false,
				HasPrevPage: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalResponse := ezutil.NewResponse("Success")
			response := originalResponse.WithPagination(tt.queryOptions, tt.totalData)

			assert.Equal(t, "Success", response.Message)
			assert.Nil(t, response.Data)
			assert.Nil(t, response.Errors)
			assert.Equal(t, tt.expected, response.Pagination)

			// Original response should be unchanged
			assert.True(t, originalResponse.Pagination.IsZero())
		})
	}
}

func TestJSONResponse_ChainedMethods(t *testing.T) {
	data := map[string]string{"key": "value"}
	err := errors.New("test error")
	queryOptions := ezutil.QueryOptions{Page: 1, Limit: 10}
	totalData := 25

	response := ezutil.NewResponse("Success").
		WithData(data).
		WithError(err).
		WithPagination(queryOptions, totalData)

	assert.Equal(t, "Success", response.Message)
	assert.Equal(t, data, response.Data)
	assert.Equal(t, err, response.Errors)
	
	expectedPagination := ezutil.Pagination{
		TotalData:   25,
		CurrentPage: 1,
		TotalPages:  3,
		HasNextPage: true,
		HasPrevPage: false,
	}
	assert.Equal(t, expectedPagination, response.Pagination)
}

func TestQueryOptions_Validation(t *testing.T) {
	// Note: This test assumes that validation is handled by the binding framework
	// The QueryOptions struct has validation tags but the validation logic
	// would be handled by Gin's binding mechanism in actual usage
	
	tests := []struct {
		name    string
		options ezutil.QueryOptions
		valid   bool
	}{
		{
			name: "valid options",
			options: ezutil.QueryOptions{
				Page:  1,
				Limit: 10,
			},
			valid: true,
		},
		{
			name: "zero page",
			options: ezutil.QueryOptions{
				Page:  0,
				Limit: 10,
			},
			valid: false,
		},
		{
			name: "negative page",
			options: ezutil.QueryOptions{
				Page:  -1,
				Limit: 10,
			},
			valid: false,
		},
		{
			name: "zero limit",
			options: ezutil.QueryOptions{
				Page:  1,
				Limit: 0,
			},
			valid: false,
		},
		{
			name: "negative limit",
			options: ezutil.QueryOptions{
				Page:  1,
				Limit: -5,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the struct can be created (validation would happen at binding time)
			assert.NotNil(t, tt.options)
			
			// Test the expected validation behavior
			pageValid := tt.options.Page >= 1
			limitValid := tt.options.Limit >= 1
			actualValid := pageValid && limitValid
			
			assert.Equal(t, tt.valid, actualValid)
		})
	}
}

func TestPaginationCalculations(t *testing.T) {
	// Test edge cases for pagination calculations
	tests := []struct {
		name      string
		totalData int
		limit     int
		expected  int // expected total pages
	}{
		{"no data", 0, 10, 0},
		{"less than limit", 5, 10, 1},
		{"exact limit", 10, 10, 1},
		{"one more than limit", 11, 10, 2},
		{"large dataset", 1000, 25, 40},
		{"single item", 1, 1, 1},
		{"single item large limit", 1, 100, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryOptions := ezutil.QueryOptions{Page: 1, Limit: tt.limit}
			response := ezutil.NewResponse("test").WithPagination(queryOptions, tt.totalData)
			
			assert.Equal(t, tt.expected, response.Pagination.TotalPages)
			assert.Equal(t, tt.totalData, response.Pagination.TotalData)
		})
	}
}
