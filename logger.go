package ezutil

import (
	"context"

	"github.com/pressly/goose/v3"
)

type Logger interface {
	Debug(args ...any)
	Info(args ...any)
	Warn(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)

	WithError(err error) Logger
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithContext(ctx context.Context) Logger

	goose.Logger
}
