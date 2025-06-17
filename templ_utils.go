package ezutil

import (
	"fmt"

	"github.com/a-h/templ"
)

func GetTemplSafeUrl(format string, args ...any) templ.SafeURL {
	return templ.URL(fmt.Sprintf(format, args...))
}
