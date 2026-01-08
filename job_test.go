package ezutil_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
)

// MockLogger for testing
type MockLogger struct {
	InfoCalls    []string
	ErrorCalls   []string
	FatalCalls   []string
	InfofCalls   []string
	FatalfCalls  []string
	ShouldPanic  bool
}

func (m *MockLogger) Debug(args ...any)                 {}
func (m *MockLogger) Info(args ...any)                  { 
	if len(args) == 0 {
		m.InfoCalls = append(m.InfoCalls, "")
	} else {
		m.InfoCalls = append(m.InfoCalls, fmt.Sprint(args...))
	}
}
func (m *MockLogger) Warn(args ...any)                  {}
func (m *MockLogger) Error(args ...any)                 { 
	if len(args) == 0 {
		m.ErrorCalls = append(m.ErrorCalls, "")
	} else {
		m.ErrorCalls = append(m.ErrorCalls, fmt.Sprint(args...))
	}
}
func (m *MockLogger) Fatal(args ...any)                 { 
	if len(args) == 0 {
		m.FatalCalls = append(m.FatalCalls, "")
	} else {
		m.FatalCalls = append(m.FatalCalls, fmt.Sprint(args...))
	}
}
func (m *MockLogger) Debugf(format string, args ...any) {}
func (m *MockLogger) Infof(format string, args ...any)  { m.InfofCalls = append(m.InfofCalls, fmt.Sprintf(format, args...)) }
func (m *MockLogger) Warnf(format string, args ...any)  {}
func (m *MockLogger) Errorf(format string, args ...any) {}
func (m *MockLogger) Fatalf(format string, args ...any) { m.FatalfCalls = append(m.FatalfCalls, fmt.Sprintf(format, args...)) }
func (m *MockLogger) Printf(format string, v ...interface{}) {}

func TestNewJob(t *testing.T) {
	logger := &MockLogger{}
	runFunc := func() error { return nil }

	job := ezutil.NewJob(logger, runFunc)
	assert.NotNil(t, job)
}

func TestNewJob_NilLogger(t *testing.T) {
	runFunc := func() error { return nil }

	assert.Panics(t, func() {
		ezutil.NewJob(nil, runFunc)
	})
}

func TestNewJob_NilRunFunc(t *testing.T) {
	logger := &MockLogger{}

	assert.Panics(t, func() {
		ezutil.NewJob(logger, nil)
	})
}

func TestJob_WithSetupFunc(t *testing.T) {
	logger := &MockLogger{}
	runFunc := func() error { return nil }
	setupFunc := func() error { return nil }

	job := ezutil.NewJob(logger, runFunc).WithSetupFunc(setupFunc)
	assert.NotNil(t, job)
}

func TestJob_WithCleanupFunc(t *testing.T) {
	logger := &MockLogger{}
	runFunc := func() error { return nil }
	cleanupFunc := func() error { return nil }

	job := ezutil.NewJob(logger, runFunc).WithCleanupFunc(cleanupFunc)
	assert.NotNil(t, job)
}

func TestJob_Run_Success(t *testing.T) {
	logger := &MockLogger{}
	runCalled := false
	runFunc := func() error {
		runCalled = true
		return nil
	}

	job := ezutil.NewJob(logger, runFunc)
	job.Run()

	assert.True(t, runCalled)
	assert.Contains(t, logger.InfoCalls, "running job...")
	assert.Len(t, logger.InfofCalls, 1)
	assert.Contains(t, logger.InfofCalls[0], "success running job for")
}

func TestJob_Run_WithSetup(t *testing.T) {
	logger := &MockLogger{}
	setupCalled := false
	runCalled := false

	setupFunc := func() error {
		setupCalled = true
		return nil
	}
	runFunc := func() error {
		runCalled = true
		return nil
	}

	job := ezutil.NewJob(logger, runFunc).WithSetupFunc(setupFunc)
	job.Run()

	assert.True(t, setupCalled)
	assert.True(t, runCalled)
	assert.Contains(t, logger.InfoCalls, "setting up job...")
	assert.Contains(t, logger.InfoCalls, "running job...")
}

func TestJob_Run_WithCleanup(t *testing.T) {
	logger := &MockLogger{}
	cleanupCalled := false
	runCalled := false

	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}
	runFunc := func() error {
		runCalled = true
		return nil
	}

	job := ezutil.NewJob(logger, runFunc).WithCleanupFunc(cleanupFunc)
	job.Run()

	assert.True(t, cleanupCalled)
	assert.True(t, runCalled)
	assert.Contains(t, logger.InfoCalls, "cleaning up job...")
	assert.Contains(t, logger.InfoCalls, "running job...")
}

func TestJob_Run_SetupError(t *testing.T) {
	logger := &MockLogger{}
	setupError := errors.New("setup failed")

	setupFunc := func() error {
		return setupError
	}
	runFunc := func() error {
		return nil
	}

	job := ezutil.NewJob(logger, runFunc).WithSetupFunc(setupFunc)
	job.Run()

	// The run function should not be called due to setup error
	// Note: In real usage, Fatal would exit, but our mock doesn't
	assert.Contains(t, logger.FatalfCalls, "error setting up job: setup failed")
}

func TestJob_Run_RunError(t *testing.T) {
	logger := &MockLogger{}
	runError := errors.New("run failed")

	runFunc := func() error {
		return runError
	}

	job := ezutil.NewJob(logger, runFunc)
	job.Run()

	assert.Contains(t, logger.FatalfCalls, "error running job: run failed")
}

func TestJob_Run_CleanupError(t *testing.T) {
	logger := &MockLogger{}
	cleanupError := errors.New("cleanup failed")

	cleanupFunc := func() error {
		return cleanupError
	}
	runFunc := func() error {
		return nil
	}

	job := ezutil.NewJob(logger, runFunc).WithCleanupFunc(cleanupFunc)
	job.Run()

	assert.Contains(t, logger.FatalfCalls, "error cleaning up job: cleanup failed")
}
