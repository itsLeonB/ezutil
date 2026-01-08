package internal_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/itsLeonB/ezutil/v2/internal"
	"github.com/stretchr/testify/assert"
)

var expectFmtOutputMsg = "Expected formatted output"

// captureOutput captures stdout for testing
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestSimpleLoggerDebug(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Debug("debug", "message", 123)
	})

	assert.Contains(t, output, "[TEST DEBUG]", "Expected output to contain '[TEST DEBUG]'")
	assert.Contains(t, output, "debug message 123", "Expected output to contain combined args")
}

func TestSimpleLoggerInfo(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Info("info message")
	})

	assert.Contains(t, output, "[TEST INFO]", "Expected output to contain '[TEST INFO]'")
	assert.Contains(t, output, "info message", "Expected output to contain 'info message'")
}

func TestSimpleLoggerWarn(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Warn("warning message")
	})

	assert.Contains(t, output, "[TEST WARN]", "Expected output to contain '[TEST WARN]'")
	assert.Contains(t, output, "warning message", "Expected output to contain 'warning message'")
}

func TestSimpleLoggerError(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Error("error message")
	})

	assert.Contains(t, output, "[TEST ERROR]", "Expected output to contain '[TEST ERROR]'")
	assert.Contains(t, output, "error message", "Expected output to contain 'error message'")
}

func TestSimpleLoggerDebugf(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Debugf("User %s has %d points", "Alice", 100)
	})

	assert.Contains(t, output, "[TEST DEBUG]", "Expected output to contain '[TEST DEBUG]'")
	assert.Contains(t, output, "User Alice has 100 points", expectFmtOutputMsg)
}

func TestSimpleLoggerInfof(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Infof("Processing %d items", 42)
	})

	assert.Contains(t, output, "[TEST INFO]", "Expected output to contain '[TEST INFO]'")
	assert.Contains(t, output, "Processing 42 items", expectFmtOutputMsg)
}

func TestSimpleLoggerWarnf(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Warnf("Temperature is %d°C", 35)
	})

	assert.Contains(t, output, "[TEST WARN]", "Expected output to contain '[TEST WARN]'")
	assert.Contains(t, output, "Temperature is 35°C", expectFmtOutputMsg)
}

func TestSimpleLoggerErrorf(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Errorf("Failed to connect to %s:%d", "localhost", 8080)
	})

	assert.Contains(t, output, "[TEST ERROR]", "Expected output to contain '[TEST ERROR]'")
	assert.Contains(t, output, "Failed to connect to localhost:8080", expectFmtOutputMsg)
}

func TestSimpleLoggerOutputFormatWithTimestamp(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Info("test message")
	})

	// Check timestamp format (HH:MM:SS.mmm)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.NotEmpty(t, lines, "Expected at least one line of output")

	// The timestamp should be at the beginning of the line
	line := lines[0]
	assert.GreaterOrEqual(t, len(line), 12, "Output line too short to contain timestamp")

	// Basic format check for timestamp (HH:MM:SS.mmm)
	timestampPart := line[:12]
	assert.Contains(t, timestampPart, ":", "Timestamp should contain colon")
	assert.Contains(t, timestampPart, ".", "Timestamp should contain dot for milliseconds")
}

func TestSimpleLoggerDifferentNamespaces(t *testing.T) {
	namespaces := []string{"API", "DB", "AUTH", "CACHE"}

	for _, ns := range namespaces {
		t.Run(ns, func(t *testing.T) {
			logger := &internal.SimpleLogger{
				Namespace: ns,
				UseColor:  false,
				MinLevel:  0,
			}

			output := captureOutput(func() {
				logger.Info("test message")
			})

			expected := fmt.Sprintf("[%s INFO]", ns)
			assert.Contains(t, output, expected, "Expected output to contain namespace")
		})
	}
}

func TestSimpleLoggerWithColor(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  true,
		MinLevel:  0,
	}

	output := captureOutput(func() {
		logger.Info("colored message")
	})

	assert.Contains(t, output, "[TEST INFO]")
	assert.Contains(t, output, "colored message")
	// Should contain ANSI color codes
	assert.Contains(t, output, "\033[36m") // Cyan color for info
	assert.Contains(t, output, "\033[0m")  // Reset color
}

func TestSimpleLoggerMinLevel(t *testing.T) {
	t.Run("debug filtered out", func(t *testing.T) {
		logger := &internal.SimpleLogger{
			Namespace: "TEST",
			UseColor:  false,
			MinLevel:  1, // Info and above
		}

		output := captureOutput(func() {
			logger.Debug("debug message")
		})

		assert.Empty(t, output, "Debug message should be filtered out")
	})

	t.Run("info allowed", func(t *testing.T) {
		logger := &internal.SimpleLogger{
			Namespace: "TEST",
			UseColor:  false,
			MinLevel:  1, // Info and above
		}

		output := captureOutput(func() {
			logger.Info("info message")
		})

		assert.Contains(t, output, "info message")
	})

	t.Run("warn level filtering", func(t *testing.T) {
		logger := &internal.SimpleLogger{
			Namespace: "TEST",
			UseColor:  false,
			MinLevel:  2, // Warn and above
		}

		output := captureOutput(func() {
			logger.Debug("debug")
			logger.Info("info")
			logger.Warn("warn")
		})

		assert.NotContains(t, output, "debug")
		assert.NotContains(t, output, "info")
		assert.Contains(t, output, "warn")
	})
}

func TestSimpleLoggerGooseInterface(t *testing.T) {
	logger := &internal.SimpleLogger{
		Namespace: "TEST",
		UseColor:  false,
		MinLevel:  0,
	}

	t.Run("Print", func(t *testing.T) {
		output := captureOutput(func() {
			logger.Print("print message")
		})

		assert.Contains(t, output, "[TEST INFO]")
		assert.Contains(t, output, "print message")
	})

	t.Run("Println", func(t *testing.T) {
		output := captureOutput(func() {
			logger.Println("println message")
		})

		assert.Contains(t, output, "[TEST INFO]")
		assert.Contains(t, output, "println message")
	})

	t.Run("Printf", func(t *testing.T) {
		output := captureOutput(func() {
			logger.Printf("printf %s %d", "message", 42)
		})

		assert.Contains(t, output, "[TEST INFO]")
		assert.Contains(t, output, "printf message 42")
	})
}

func TestSimpleLoggerColorCodes(t *testing.T) {
	tests := []struct {
		level    string
		logFunc  func(*internal.SimpleLogger)
		expected string
	}{
		{"DEBUG", func(l *internal.SimpleLogger) { l.Debug("test") }, "\033[36m"},
		{"INFO", func(l *internal.SimpleLogger) { l.Info("test") }, "\033[36m"},
		{"WARN", func(l *internal.SimpleLogger) { l.Warn("test") }, "\033[33m"},
		{"ERROR", func(l *internal.SimpleLogger) { l.Error("test") }, "\033[31m"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger := &internal.SimpleLogger{
				Namespace: "TEST",
				UseColor:  true,
				MinLevel:  0,
			}

			output := captureOutput(func() {
				tt.logFunc(logger)
			})

			assert.Contains(t, output, tt.expected, "Should contain correct color code")
			assert.Contains(t, output, "\033[0m", "Should contain reset code")
		})
	}
}

// Note: Testing Fatal and Fatalf methods that call os.Exit(1) would terminate the test process.
// In a real-world scenario, you might want to refactor the logger to accept an interface
// for exiting (dependency injection) to make it testable, or use integration tests
// that run the code in a separate process.

// Benchmark tests
func BenchmarkSimpleLoggerInfo(b *testing.B) {
	logger := &internal.SimpleLogger{
		Namespace: "BENCH",
		UseColor:  false,
		MinLevel:  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkSimpleLoggerInfof(b *testing.B) {
	logger := &internal.SimpleLogger{
		Namespace: "BENCH",
		UseColor:  false,
		MinLevel:  1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Infof("benchmark message %d", i)
	}
}

func BenchmarkSimpleLoggerFilteredOut(b *testing.B) {
	logger := &internal.SimpleLogger{
		Namespace: "BENCH",
		UseColor:  false,
		MinLevel:  2, // Warn and above
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("this should be filtered out")
	}
}
