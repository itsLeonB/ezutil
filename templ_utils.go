package ezutil

import (
	"fmt"

	"github.com/a-h/templ"
)

// GetTemplSafeUrl creates a templ.SafeURL from a format string and arguments.
// It uses fmt.Sprintf to format the URL and wraps it in templ.URL for safe template rendering.
// This prevents XSS attacks by ensuring URLs are properly escaped in templates.
func GetTemplSafeUrl(format string, args ...any) templ.SafeURL {
	return templ.URL(fmt.Sprintf(format, args...))
}
