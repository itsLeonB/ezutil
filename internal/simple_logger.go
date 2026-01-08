// Package internal provides internal utilities for the ezutil package.
package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type logLevel string

const (
	levelDebug logLevel = "DEBUG"
	levelInfo  logLevel = "INFO"
	levelWarn  logLevel = "WARN"
	levelError logLevel = "ERROR"
	levelFatal logLevel = "FATAL"
)

type SimpleLogger struct {
	Namespace string
	UseColor  bool
	MinLevel  int
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

func (s *SimpleLogger) output(level logLevel, msg string) {
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

func (s *SimpleLogger) outputf(level logLevel, format string, args ...any) {
	if levelToInt[level] < s.MinLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	s.output(level, msg)
}

// Debug logs a debug message.
func (s *SimpleLogger) Debug(args ...any) {
	s.output(levelDebug, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Info logs an info message.
func (s *SimpleLogger) Info(args ...any) {
	s.output(levelInfo, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Warn logs a warning message.
func (s *SimpleLogger) Warn(args ...any) {
	s.output(levelWarn, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Error logs an error message.
func (s *SimpleLogger) Error(args ...any) {
	s.output(levelError, strings.TrimRight(fmt.Sprintln(args...), "\n"))
}

// Fatal logs a fatal message and exits the program.
func (s *SimpleLogger) Fatal(args ...any) {
	s.output(levelFatal, strings.TrimRight(fmt.Sprintln(args...), "\n"))
	os.Exit(1)
}

// Debugf logs a formatted debug message.
func (s *SimpleLogger) Debugf(format string, args ...any) { s.outputf(levelDebug, format, args...) }

// Infof logs a formatted info message.
func (s *SimpleLogger) Infof(format string, args ...any)  { s.outputf(levelInfo, format, args...) }

// Warnf logs a formatted warning message.
func (s *SimpleLogger) Warnf(format string, args ...any)  { s.outputf(levelWarn, format, args...) }

// Errorf logs a formatted error message.
func (s *SimpleLogger) Errorf(format string, args ...any) { s.outputf(levelError, format, args...) }

// Fatalf logs a formatted fatal message and exits the program.
func (s *SimpleLogger) Fatalf(format string, args ...any) {
	s.outputf(levelFatal, format, args...)
	os.Exit(1)
}

// Print logs a message using Info level (goose.Logger interface).
func (s *SimpleLogger) Print(args ...any)                 { s.Info(fmt.Sprint(args...)) }

// Println logs a message using Info level (goose.Logger interface).
func (s *SimpleLogger) Println(args ...any)               { s.Info(strings.TrimRight(fmt.Sprintln(args...), "\n")) }

// Printf logs a formatted message using Info level (goose.Logger interface).
func (s *SimpleLogger) Printf(format string, args ...any) { s.Infof(format, args...) }
