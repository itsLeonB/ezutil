package ezutil_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupGinTest() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestGetPathParam(t *testing.T) {
	tests := []struct {
		name        string
		paramKey    string
		paramValue  string
		expectValue interface{}
		expectExist bool
		expectError bool
		paramType   string
	}{
		{
			name:        "string parameter exists",
			paramKey:    "name",
			paramValue:  "john",
			expectValue: "john",
			expectExist: true,
			expectError: false,
			paramType:   "string",
		},
		{
			name:        "int parameter exists",
			paramKey:    "id",
			paramValue:  "123",
			expectValue: 123,
			expectExist: true,
			expectError: false,
			paramType:   "int",
		},
		{
			name:        "bool parameter true",
			paramKey:    "active",
			paramValue:  "true",
			expectValue: true,
			expectExist: true,
			expectError: false,
			paramType:   "bool",
		},
		{
			name:        "bool parameter false",
			paramKey:    "active",
			paramValue:  "false",
			expectValue: false,
			expectExist: true,
			expectError: false,
			paramType:   "bool",
		},
		{
			name:        "uuid parameter",
			paramKey:    "uuid",
			paramValue:  "550e8400-e29b-41d4-a716-446655440000",
			expectValue: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			expectExist: true,
			expectError: false,
			paramType:   "uuid",
		},
		{
			name:        "parameter does not exist",
			paramKey:    "missing",
			paramValue:  "",
			expectValue: "",
			expectExist: false,
			expectError: false,
			paramType:   "string",
		},
		{
			name:        "invalid int parameter",
			paramKey:    "id",
			paramValue:  "abc",
			expectValue: 0,
			expectExist: false, // This test has complex logic, just check error
			expectError: true,
			paramType:   "int",
		},
		{
			name:        "invalid bool parameter",
			paramKey:    "active",
			paramValue:  "maybe",
			expectValue: false,
			expectExist: false, // This test has complex logic, just check error
			expectError: true,
			paramType:   "bool",
		},
		{
			name:        "invalid uuid parameter",
			paramKey:    "uuid",
			paramValue:  "invalid-uuid",
			expectValue: uuid.UUID{},
			expectExist: false, // This test has complex logic, just check error
			expectError: true,
			paramType:   "uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupGinTest()
			
			// Set up the parameter if it should exist or if we expect an error (which means param exists but is invalid)
			if tt.expectExist || tt.expectError {
				c.Params = gin.Params{
					{Key: tt.paramKey, Value: tt.paramValue},
				}
			}

			switch tt.paramType {
			case "string":
				value, exists, err := ezutil.GetPathParam[string](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
					// For error cases, just check that error occurred, don't check exists value
					// since the behavior depends on implementation details
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
					assert.Equal(t, tt.expectExist, exists)
				}

			case "int":
				value, exists, err := ezutil.GetPathParam[int](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
					// For error cases, just check that error occurred
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
					assert.Equal(t, tt.expectExist, exists)
				}

			case "bool":
				value, exists, err := ezutil.GetPathParam[bool](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
					// For error cases, just check that error occurred
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
					assert.Equal(t, tt.expectExist, exists)
				}

			case "uuid":
				value, exists, err := ezutil.GetPathParam[uuid.UUID](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
					// For error cases, just check that error occurred
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
					assert.Equal(t, tt.expectExist, exists)
				}
			}
		})
	}
}

func TestGetRequiredPathParam(t *testing.T) {
	tests := []struct {
		name        string
		paramKey    string
		paramValue  string
		expectValue interface{}
		expectError bool
		paramType   string
	}{
		{
			name:        "string parameter exists",
			paramKey:    "name",
			paramValue:  "john",
			expectValue: "john",
			expectError: false,
			paramType:   "string",
		},
		{
			name:        "int parameter exists",
			paramKey:    "id",
			paramValue:  "123",
			expectValue: 123,
			expectError: false,
			paramType:   "int",
		},
		{
			name:        "parameter missing",
			paramKey:    "missing",
			paramValue:  "",
			expectValue: nil,
			expectError: true,
			paramType:   "string",
		},
		{
			name:        "invalid parameter",
			paramKey:    "id",
			paramValue:  "abc",
			expectValue: 0,
			expectError: true,
			paramType:   "int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupGinTest()
			
			// Set up the parameter if it should exist
			if tt.paramValue != "" {
				c.Params = gin.Params{
					{Key: tt.paramKey, Value: tt.paramValue},
				}
			}

			switch tt.paramType {
			case "string":
				value, err := ezutil.GetRequiredPathParam[string](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
				}

			case "int":
				value, err := ezutil.GetRequiredPathParam[int](c, tt.paramKey)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
				}
			}
		})
	}
}

func TestBindRequest(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" form:"name" binding:"required"`
		Age   int    `json:"age" form:"age" binding:"required,min=1"`
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	tests := []struct {
		name        string
		requestBody string
		contentType string
		bindingType binding.Binding
		expectError bool
		expected    TestStruct
	}{
		{
			name:        "valid JSON binding",
			requestBody: `{"name":"John","age":30,"email":"john@example.com"}`,
			contentType: "application/json",
			bindingType: binding.JSON,
			expectError: false,
			expected: TestStruct{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
		},
		{
			name:        "invalid JSON",
			requestBody: `{"name":"John","age":}`,
			contentType: "application/json",
			bindingType: binding.JSON,
			expectError: true,
		},
		{
			name:        "missing required field",
			requestBody: `{"name":"John"}`,
			contentType: "application/json",
			bindingType: binding.JSON,
			expectError: true,
		},
		{
			name:        "form binding",
			requestBody: "name=John&age=30&email=john@example.com",
			contentType: "application/x-www-form-urlencoded",
			bindingType: binding.Form,
			expectError: false,
			expected: TestStruct{
				Name:  "John",
				Age:   30,
				Email: "john@example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupGinTest()
			
			// Set up request
			var req *http.Request
			if tt.bindingType == binding.Form {
				req = httptest.NewRequest("POST", "/test", strings.NewReader(tt.requestBody))
				req.Header.Set("Content-Type", tt.contentType)
			} else {
				req = httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.requestBody))
				req.Header.Set("Content-Type", tt.contentType)
			}
			c.Request = req

			result, err := ezutil.BindRequest[TestStruct](c, tt.bindingType)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetFromContext(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       interface{}
		expectError bool
		valueType   string
	}{
		{
			name:        "string value exists",
			key:         "user_id",
			value:       "123",
			expectError: false,
			valueType:   "string",
		},
		{
			name:        "int value exists",
			key:         "count",
			value:       42,
			expectError: false,
			valueType:   "int",
		},
		{
			name:        "bool value exists",
			key:         "active",
			value:       true,
			expectError: false,
			valueType:   "bool",
		},
		{
			name:        "key does not exist",
			key:         "missing",
			value:       nil,
			expectError: true,
			valueType:   "string",
		},
		{
			name:        "wrong type assertion",
			key:         "user_id",
			value:       "123", // string value
			expectError: true,
			valueType:   "int", // but we expect int
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupGinTest()
			
			// Set up context value if it should exist
			if tt.value != nil && !tt.expectError {
				c.Set(tt.key, tt.value)
			} else if tt.name == "wrong type assertion" {
				c.Set(tt.key, tt.value)
			}

			switch tt.valueType {
			case "string":
				value, err := ezutil.GetFromContext[string](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.value, value)
				}

			case "int":
				value, err := ezutil.GetFromContext[int](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.value, value)
				}

			case "bool":
				value, err := ezutil.GetFromContext[bool](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.value, value)
				}
			}
		})
	}
}

func TestGetAndParseFromContext(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       string
		expectValue interface{}
		expectError bool
		parseType   string
	}{
		{
			name:        "parse string to int",
			key:         "user_id",
			value:       "123",
			expectValue: 123,
			expectError: false,
			parseType:   "int",
		},
		{
			name:        "parse string to bool",
			key:         "active",
			value:       "true",
			expectValue: true,
			expectError: false,
			parseType:   "bool",
		},
		{
			name:        "parse string to uuid",
			key:         "uuid",
			value:       "550e8400-e29b-41d4-a716-446655440000",
			expectValue: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			expectError: false,
			parseType:   "uuid",
		},
		{
			name:        "key does not exist",
			key:         "missing",
			value:       "",
			expectValue: nil,
			expectError: true,
			parseType:   "int",
		},
		{
			name:        "invalid parse",
			key:         "user_id",
			value:       "abc",
			expectValue: 0,
			expectError: true,
			parseType:   "int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := setupGinTest()
			
			// Set up context value if it should exist
			if tt.value != "" {
				c.Set(tt.key, tt.value)
			}

			switch tt.parseType {
			case "int":
				value, err := ezutil.GetAndParseFromContext[int](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
				}

			case "bool":
				value, err := ezutil.GetAndParseFromContext[bool](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
				}

			case "uuid":
				value, err := ezutil.GetAndParseFromContext[uuid.UUID](c, tt.key)
				if tt.expectError {
					assert.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.expectValue, value)
				}
			}
		})
	}
}

func TestGinUtilsIntegration(t *testing.T) {
	// Test a realistic scenario combining multiple utilities
	t.Run("user profile endpoint simulation", func(t *testing.T) {
		c, _ := setupGinTest()
		
		// Set up path parameters
		c.Params = gin.Params{
			{Key: "user_id", Value: "123"},
			{Key: "profile_id", Value: "456"},
		}
		
		// Set up context values
		c.Set("role", "admin")
		c.Set("permissions", "read,write")
		
		// Set up request body
		requestBody := `{"name":"John Doe","email":"john@example.com"}`
		req := httptest.NewRequest("PUT", "/users/123/profiles/456", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		// Test path parameter extraction
		userID, exists, err := ezutil.GetPathParam[int](c, "user_id")
		require.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, 123, userID)

		profileID, err := ezutil.GetRequiredPathParam[int](c, "profile_id")
		require.NoError(t, err)
		assert.Equal(t, 456, profileID)

		// Test context value retrieval
		role, err := ezutil.GetFromContext[string](c, "role")
		require.NoError(t, err)
		assert.Equal(t, "admin", role)

		// Test request binding
		type UpdateRequest struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		updateReq, err := ezutil.BindRequest[UpdateRequest](c, binding.JSON)
		require.NoError(t, err)
		assert.Equal(t, "John Doe", updateReq.Name)
		assert.Equal(t, "john@example.com", updateReq.Email)
	})
}
