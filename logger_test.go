package ezutil_test

import (
	"testing"

	"github.com/itsLeonB/ezutil/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleLogger(t *testing.T) {
	logger := ezutil.NewSimpleLogger("TEST", false, 0)
	assert.NotNil(t, logger)
}

func TestNewSimpleLogger_WithColor(t *testing.T) {
	logger := ezutil.NewSimpleLogger("TEST", true, 1)
	assert.NotNil(t, logger)
}

func TestLogger_Interface(t *testing.T) {
	logger := ezutil.NewSimpleLogger("TEST", false, 0)
	
	// Test that it implements the Logger interface
	var _ ezutil.Logger = logger
	
	// Test basic methods don't panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	logger.Debugf("test %s", "format")
	logger.Infof("test %s", "format")
	logger.Warnf("test %s", "format")
	logger.Errorf("test %s", "format")
	logger.Printf("test %s", "format")
}
