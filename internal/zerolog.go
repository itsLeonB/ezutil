package internal

import (
	"fmt"
	"io"

	"github.com/rs/zerolog"
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
func (s *zerologAdapter) Debug(args ...any) {
	s.logger.Debug().Msg(fmt.Sprint(args...))
}

// Info logs an info message.
func (s *zerologAdapter) Info(args ...any) {
	s.logger.Info().Msg(fmt.Sprint(args...))
}

// Warn logs a warning message.
func (s *zerologAdapter) Warn(args ...any) {
	s.logger.Warn().Msg(fmt.Sprint(args...))
}

// Error logs an error message.
func (s *zerologAdapter) Error(args ...any) {
	s.logger.Error().Msg(fmt.Sprint(args...))
}

// Fatal logs a fatal message and exits the program.
func (s *zerologAdapter) Fatal(args ...any) {
	s.logger.Fatal().Msg(fmt.Sprint(args...))
}

// Debugf logs a formatted debug message.
func (s *zerologAdapter) Debugf(format string, args ...any) {
	s.logger.Debug().Msg(fmt.Sprintf(format, args...))
}

// Infof logs a formatted info message.
func (s *zerologAdapter) Infof(format string, args ...any) {
	s.logger.Info().Msg(fmt.Sprintf(format, args...))
}

// Warnf logs a formatted warning message.
func (s *zerologAdapter) Warnf(format string, args ...any) {
	s.logger.Warn().Msg(fmt.Sprintf(format, args...))
}

// Errorf logs a formatted error message.
func (s *zerologAdapter) Errorf(format string, args ...any) {
	s.logger.Error().Msg(fmt.Sprintf(format, args...))
}

// Fatalf logs a formatted fatal message and exits the program.
func (s *zerologAdapter) Fatalf(format string, args ...any) {
	s.logger.Fatal().Msg(fmt.Sprintf(format, args...))
}

// Print logs a message using Info level (goose.Logger interface).
func (s *zerologAdapter) Print(args ...any) { s.Info(args...) }

// Println logs a message using Info level (goose.Logger interface).
func (s *zerologAdapter) Println(args ...any) { s.Info(args...) }

// Printf logs a formatted message using Info level (goose.Logger interface).
func (s *zerologAdapter) Printf(format string, args ...any) { s.Infof(format, args...) }
