// Package simple provides simple implementation of Logger that logs to STDOUT
package simple

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/itsLeonB/ezutil/v2"
)

type logLevel string

const (
	levelDebug logLevel = "DEBUG"
	levelInfo  logLevel = "INFO"
	levelWarn  logLevel = "WARN"
	levelError logLevel = "ERROR"
	levelFatal logLevel = "FATAL"
)

type Logger struct {
	Namespace string
	UseColor  bool
	MinLevel  int
}

func NewLogger(namespace string, useColor bool, minLevel int) *Logger {
	return &Logger{
		Namespace: namespace,
		UseColor:  useColor,
		MinLevel:  minLevel,
	}
}

var colors = map[logLevel]string{
	levelDebug: "\033[36m",
	levelInfo:  "\033[36m",
	levelWarn:  "\033[33m",
	levelError: "\033[31m",
	levelFatal: "\033[31m",
}

var levelToInt = map[logLevel]int{
	levelDebug: 0,
	levelInfo:  1,
	levelWarn:  2,
	levelError: 3,
	levelFatal: 4,
}

func (s *Logger) output(level logLevel, msg string) {
	if levelToInt[level] < s.MinLevel {
		return
	}
	var colorStart, colorReset string
	if s.UseColor {
		colorStart = colors[level]
		colorReset = "\033[0m"
	}
	fmt.Printf("%s%s [%s %s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, s.Namespace, level, msg, colorReset)
}

func (s *Logger) outputf(level logLevel, format string, args ...any) {
	if levelToInt[level] < s.MinLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	s.output(level, msg)
}

// Debug logs a debug message.
func (s *Logger) Debug(args ...any) {
	s.output(levelDebug, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Info logs an info message.
func (s *Logger) Info(args ...any) {
	s.output(levelInfo, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Warn logs a warning message.
func (s *Logger) Warn(args ...any) {
	s.output(levelWarn, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Error logs an error message.
func (s *Logger) Error(args ...any) {
	s.output(levelError, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Fatal logs a fatal message and exits the program.
func (s *Logger) Fatal(args ...any) {
	s.output(levelFatal, strings.TrimRight(fmt.Sprintln(args...), "\n"))
	os.Exit(1)
}

// Debugf logs a formatted debug message.
func (s *Logger) Debugf(format string, args ...any) { s.outputf(levelDebug, format, args...) }

// Infof logs a formatted info message.
func (s *Logger) Infof(format string, args ...any) { s.outputf(levelInfo, format, args...) }

// Warnf logs a formatted warning message.
func (s *Logger) Warnf(format string, args ...any) { s.outputf(levelWarn, format, args...) }

// Errorf logs a formatted error message.
func (s *Logger) Errorf(format string, args ...any) { s.outputf(levelError, format, args...) }

// Fatalf logs a formatted fatal message and exits the program.
func (s *Logger) Fatalf(format string, args ...any) {
	s.outputf(levelFatal, format, args...)
	os.Exit(1)
}

// Print logs a message using Info level (goose.Logger interface).
func (s *Logger) Print(args ...any) { s.output(levelInfo, fmt.Sprint(args...)) }

// Println logs a message using Info level (goose.Logger interface).
func (s *Logger) Println(args ...any) {
	s.output(levelInfo, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Printf logs a formatted message using Info level (goose.Logger interface).
func (s *Logger) Printf(format string, args ...any) {
	s.output(levelInfo, fmt.Sprintf(format, args...))
}

func (s *Logger) WithError(err error) ezutil.Logger {
	return s
}

func (s *Logger) WithField(key string, value any) ezutil.Logger {
	return s
}

func (s *Logger) WithFields(fields map[string]any) ezutil.Logger {
	return s
}

func (s *Logger) WithContext(ctx context.Context) ezutil.Logger {
	return s
}
