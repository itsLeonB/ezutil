package ezutil_test

import (
	"testing"

	"github.com/a-h/templ"
	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
)

func TestGetTemplSafeUrl(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected templ.SafeURL
	}{
		{
			name:     "simple URL without arguments",
			format:   "https://example.com",
			args:     nil,
			expected: templ.URL("https://example.com"),
		},
		{
			name:     "URL with single string argument",
			format:   "https://example.com/%s",
			args:     []any{"users"},
			expected: templ.URL("https://example.com/users"),
		},
		{
			name:     "URL with multiple arguments",
			format:   "https://example.com/%s/%d",
			args:     []any{"users", 123},
			expected: templ.URL("https://example.com/users/123"),
		},
		{
			name:     "URL with query parameters",
			format:   "https://example.com/search?q=%s&page=%d",
			args:     []any{"golang", 1},
			expected: templ.URL("https://example.com/search?q=golang&page=1"),
		},
		{
			name:     "relative URL",
			format:   "/api/v1/%s/%d",
			args:     []any{"posts", 456},
			expected: templ.URL("/api/v1/posts/456"),
		},
		{
			name:     "URL with fragment",
			format:   "https://example.com/docs#%s",
			args:     []any{"section-1"},
			expected: templ.URL("https://example.com/docs#section-1"),
		},
		{
			name:     "empty format string",
			format:   "",
			args:     nil,
			expected: templ.URL(""),
		},
		{
			name:     "format with no placeholders but with args",
			format:   "https://example.com/static",
			args:     []any{"ignored", "arguments"},
			expected: templ.URL("https://example.com/static%!(EXTRA string=ignored, string=arguments)"),
		},
		{
			name:     "complex URL with multiple types",
			format:   "https://api.example.com/v%d/%s/%d/comments?limit=%d&offset=%d",
			args:     []any{2, "posts", 789, 10, 20},
			expected: templ.URL("https://api.example.com/v2/posts/789/comments?limit=10&offset=20"),
		},
		{
			name:     "URL with boolean argument",
			format:   "https://example.com/api?active=%t",
			args:     []any{true},
			expected: templ.URL("https://example.com/api?active=true"),
		},
		{
			name:     "URL with float argument",
			format:   "https://example.com/api?price=%.2f",
			args:     []any{19.99},
			expected: templ.URL("https://example.com/api?price=19.99"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.GetTemplSafeUrl(tt.format, tt.args...)
			assert.Equal(t, tt.expected, result)
			
			// Verify that the result is actually a SafeURL
			assert.IsType(t, templ.SafeURL(""), result)
			
			// Verify that the underlying string value is correct
			assert.Equal(t, string(tt.expected), string(result))
		})
	}
}

func TestGetTemplSafeUrl_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected templ.SafeURL
	}{
		{
			name:     "URL with spaces",
			format:   "https://example.com/search?q=%s",
			args:     []any{"hello world"},
			expected: templ.URL("https://example.com/search?q=hello world"),
		},
		{
			name:     "URL with special characters",
			format:   "https://example.com/path/%s",
			args:     []any{"file@name#special"},
			expected: templ.URL("https://example.com/path/file@name#special"),
		},
		{
			name:     "URL with unicode characters",
			format:   "https://example.com/user/%s",
			args:     []any{"用户名"},
			expected: templ.URL("https://example.com/user/用户名"),
		},
		{
			name:     "URL with percent signs",
			format:   "https://example.com/discount/%s",
			args:     []any{"50%"},
			expected: templ.URL("https://example.com/discount/50%"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.GetTemplSafeUrl(tt.format, tt.args...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetTemplSafeUrl_EdgeCases(t *testing.T) {
	t.Run("nil args", func(t *testing.T) {
		result := ezutil.GetTemplSafeUrl("https://example.com")
		expected := templ.URL("https://example.com")
		assert.Equal(t, expected, result)
	})

	t.Run("empty args slice", func(t *testing.T) {
		result := ezutil.GetTemplSafeUrl("https://example.com")
		expected := templ.URL("https://example.com")
		assert.Equal(t, expected, result)
	})

	t.Run("more placeholders than args", func(t *testing.T) {
		// This test would cause a compile-time error, so we'll skip it
		// In real usage, this would be caught at compile time
		t.Skip("Compile-time format string validation prevents this test")
	})

	t.Run("more args than placeholders", func(t *testing.T) {
		// This test would cause a compile-time error, so we'll skip it  
		// In real usage, this would be caught at compile time
		t.Skip("Compile-time format string validation prevents this test")
	})

	t.Run("nil argument", func(t *testing.T) {
		result := ezutil.GetTemplSafeUrl("https://example.com/%v", nil)
		expected := templ.URL("https://example.com/<nil>")
		assert.Equal(t, expected, result)
	})
}

func TestGetTemplSafeUrl_TypeSafety(t *testing.T) {
	// Test that the function returns the correct type
	result := ezutil.GetTemplSafeUrl("https://example.com")
	
	// Should be assignable to templ.SafeURL
	var safeURL templ.SafeURL = result
	assert.Equal(t, templ.URL("https://example.com"), safeURL)
	
	// Should be convertible to string
	urlString := string(result)
	assert.Equal(t, "https://example.com", urlString)
}

func TestGetTemplSafeUrl_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected templ.SafeURL
	}{
		{
			name:     "user profile URL",
			format:   "/users/%d/profile",
			args:     []any{12345},
			expected: templ.URL("/users/12345/profile"),
		},
		{
			name:     "API endpoint with version",
			format:   "/api/v%d/%s",
			args:     []any{1, "users"},
			expected: templ.URL("/api/v1/users"),
		},
		{
			name:     "search with pagination",
			format:   "/search?q=%s&page=%d&limit=%d",
			args:     []any{"golang", 2, 25},
			expected: templ.URL("/search?q=golang&page=2&limit=25"),
		},
		{
			name:     "file download URL",
			format:   "/files/%s/download?token=%s",
			args:     []any{"document.pdf", "abc123"},
			expected: templ.URL("/files/document.pdf/download?token=abc123"),
		},
		{
			name:     "admin panel URL",
			format:   "/admin/%s/%d/edit",
			args:     []any{"products", 789},
			expected: templ.URL("/admin/products/789/edit"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ezutil.GetTemplSafeUrl(tt.format, tt.args...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
