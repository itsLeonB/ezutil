// Package zerolog provides Logger implementation that uses zerolog under the hood
package zerolog

import (
	"fmt"
	"io"

	"github.com/itsLeonB/ezutil/v2/gorm"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
)

type zerologAdapter struct {
	logger zerolog.Logger
}

func NewZerologAdapter(namespace string, writer io.Writer) *zerologAdapter {
	return &zerologAdapter{
		zerolog.New(writer).With().Timestamp().Logger(),
	}
}

// Debug logs a debug message.
func (z *zerologAdapter) Debug(args ...any) {
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

// Info logs an info message.
func (z *zerologAdapter) Info(args ...any) {
	z.logger.Info().Msg(fmt.Sprint(args...))
}

// Warn logs a warning message.
func (z *zerologAdapter) Warn(args ...any) {
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

// Error logs an error message.
func (z *zerologAdapter) Error(args ...any) {
	z.logger.Error().Msg(fmt.Sprint(args...))
}

// Fatal logs a fatal message and exits the program.
func (z *zerologAdapter) Fatal(args ...any) {
	z.logger.Fatal().Msg(fmt.Sprint(args...))
}

// Debugf logs a formatted debug message.
func (z *zerologAdapter) Debugf(format string, args ...any) {
	z.logger.Debug().Msg(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message.
func (z *zerologAdapter) Infof(format string, args ...any) {
	z.logger.Info().Msg(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message.
func (z *zerologAdapter) Warnf(format string, args ...any) {
	z.logger.Warn().Msg(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message.
func (z *zerologAdapter) Errorf(format string, args ...any) {
	z.logger.Error().Msg(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal message and exits the program.
func (z *zerologAdapter) Fatalf(format string, args ...any) {
	z.logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

// Print logs a message using Info level (goose.Logger interface).
func (z *zerologAdapter) Print(args ...any) { z.Info(args...) }

// Println logs a message using Info level (goose.Logger interface).
func (z *zerologAdapter) Println(args ...any) { z.Info(args...) }

// Printf logs a formatted message using Info level (goose.Logger interface).
func (z *zerologAdapter) Printf(format string, args ...any) { z.Infof(format, args...) }

func (z *zerologAdapter) Zerolog() zerolog.Logger {
	return z.logger
}

func (z *zerologAdapter) AsGorm() logger.Interface {
	return gorm.NewGormLogger(z)
}
