package ezutil

import (
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
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

	Zerolog() zerolog.Logger
	AsGorm() logger.Interface

	goose.Logger
}
