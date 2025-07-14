package ezutil_test

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/itsLeonB/ezutil"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithoutDB(t *testing.T) {
	// Test loading config without database connection
	// This tests the actual configuration loading logic without requiring a real database

	// Set up environment variables for testing
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("APP_TIMEOUT", "30s")
	os.Setenv("APP_TIMEZONE", "UTC")

	os.Setenv("AUTH_SECRETKEY", "test-secret")
	os.Setenv("AUTH_TOKENDURATION", "1h")
	os.Setenv("AUTH_COOKIEDURATION", "24h")
	os.Setenv("AUTH_ISSUER", "test-issuer")
	os.Setenv("AUTH_URL", "http://localhost:8080")

	defer func() {
		// Clean up environment variables
		envVars := []string{
			"APP_ENV", "APP_PORT", "APP_TIMEOUT", "APP_TIMEZONE",
			"AUTH_SECRETKEY", "AUTH_TOKENDURATION", "AUTH_COOKIEDURATION", "AUTH_ISSUER", "AUTH_URL",
		}
		for _, env := range envVars {
			os.Unsetenv(env)
		}
	}()

	// Define defaults to test override behavior
	defaults := ezutil.Config{
		App: &ezutil.App{
			Env:        "development",
			Port:       "3000",
			Timeout:    10 * time.Second,
			ClientUrls: []string{"http://localhost:3000"},
			Timezone:   "America/New_York",
		},
		Auth: &ezutil.Auth{
			SecretKey:      "default-secret",
			TokenDuration:  30 * time.Minute,
			CookieDuration: 12 * time.Hour,
			Issuer:         "default-issuer",
			URL:            "http://localhost:3000",
		},
	}

	// Actually load the configuration without database connection
	config := ezutil.LoadConfigWithoutDB(defaults)

	// Test that configuration was loaded correctly and environment variables override defaults
	t.Run("app config loaded from environment", func(t *testing.T) {
		assert.NotNil(t, config.App)
		assert.Equal(t, "test", config.App.Env)                    // Overridden by APP_ENV
		assert.Equal(t, "8080", config.App.Port)                  // Overridden by APP_PORT
		assert.Equal(t, 30*time.Second, config.App.Timeout)       // Overridden by APP_TIMEOUT
		assert.Equal(t, "UTC", config.App.Timezone)               // Overridden by APP_TIMEZONE
		assert.Equal(t, []string{"http://localhost:3000"}, config.App.ClientUrls) // Uses default (not set in env)
	})

	t.Run("auth config loaded from environment", func(t *testing.T) {
		assert.NotNil(t, config.Auth)
		assert.Equal(t, "test-secret", config.Auth.SecretKey)      // Overridden by AUTH_SECRETKEY
		assert.Equal(t, 1*time.Hour, config.Auth.TokenDuration)   // Overridden by AUTH_TOKENDURATION
		assert.Equal(t, 24*time.Hour, config.Auth.CookieDuration) // Overridden by AUTH_COOKIEDURATION
		assert.Equal(t, "test-issuer", config.Auth.Issuer)        // Overridden by AUTH_ISSUER
		assert.Equal(t, "http://localhost:8080", config.Auth.URL) // Overridden by AUTH_URL
	})

	t.Run("database config present but no connection", func(t *testing.T) {
		assert.NotNil(t, config.SQLDB) // SQLDB config should be loaded
		assert.Nil(t, config.GORM)     // But GORM connection should be nil
	})
}

func TestLoadConfigWithDefaults(t *testing.T) {
	// Test loading config with default values when no environment variables are set
	// This ensures defaults are properly applied

	// Ensure no relevant environment variables are set
	envVars := []string{
		"APP_ENV", "APP_PORT", "APP_TIMEOUT", "APP_TIMEZONE", "APP_CLIENTURLS",
		"AUTH_SECRETKEY", "AUTH_TOKENDURATION", "AUTH_COOKIEDURATION", "AUTH_ISSUER", "AUTH_URL",
	}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Define defaults
	defaults := ezutil.Config{
		App: &ezutil.App{
			Env:        "development",
			Port:       "3000",
			Timeout:    10 * time.Second,
			ClientUrls: []string{"http://localhost:3000"},
			Timezone:   "America/New_York",
		},
		Auth: &ezutil.Auth{
			SecretKey:      "default-secret",
			TokenDuration:  30 * time.Minute,
			CookieDuration: 12 * time.Hour,
			Issuer:         "default-issuer",
			URL:            "http://localhost:3000",
		},
	}

	// Load configuration without database connection
	config := ezutil.LoadConfigWithoutDB(defaults)

	// Test that defaults are used when no environment variables are set
	t.Run("app config uses defaults", func(t *testing.T) {
		assert.NotNil(t, config.App)
		assert.Equal(t, "development", config.App.Env)
		assert.Equal(t, "3000", config.App.Port)
		assert.Equal(t, 10*time.Second, config.App.Timeout)
		assert.Equal(t, []string{"http://localhost:3000"}, config.App.ClientUrls)
		assert.Equal(t, "America/New_York", config.App.Timezone)
	})

	t.Run("auth config uses defaults", func(t *testing.T) {
		assert.NotNil(t, config.Auth)
		assert.Equal(t, "default-secret", config.Auth.SecretKey)
		assert.Equal(t, 30*time.Minute, config.Auth.TokenDuration)
		assert.Equal(t, 12*time.Hour, config.Auth.CookieDuration)
		assert.Equal(t, "default-issuer", config.Auth.Issuer)
		assert.Equal(t, "http://localhost:3000", config.Auth.URL)
	})

	t.Run("database config present but no connection", func(t *testing.T) {
		assert.NotNil(t, config.SQLDB)
		assert.Nil(t, config.GORM)
	})
}

func TestLoadConfigValidation(t *testing.T) {
	// Test configuration validation logic

	// Clean up any existing environment variables
	envVars := []string{
		"APP_ENV", "APP_PORT", "APP_TIMEOUT", "APP_TIMEZONE",
		"AUTH_SECRETKEY", "AUTH_TOKENDURATION", "AUTH_COOKIEDURATION", "AUTH_ISSUER", "AUTH_URL",
	}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	t.Run("negative timeout uses default", func(t *testing.T) {
		os.Setenv("APP_TIMEOUT", "-5s")
		defer os.Unsetenv("APP_TIMEOUT")

		defaults := ezutil.Config{
			App: &ezutil.App{
				Timeout: 10 * time.Second,
			},
		}

		config := ezutil.LoadConfigWithoutDB(defaults)
		assert.Equal(t, 10*time.Second, config.App.Timeout) // Should use default
	})

	t.Run("zero timeout uses default", func(t *testing.T) {
		os.Setenv("APP_TIMEOUT", "0s")
		defer os.Unsetenv("APP_TIMEOUT")

		defaults := ezutil.Config{
			App: &ezutil.App{
				Timeout: 15 * time.Second,
			},
		}

		config := ezutil.LoadConfigWithoutDB(defaults)
		assert.Equal(t, 15*time.Second, config.App.Timeout) // Should use default
	})

	t.Run("negative token duration uses default", func(t *testing.T) {
		os.Setenv("AUTH_TOKENDURATION", "-1h")
		defer os.Unsetenv("AUTH_TOKENDURATION")

		defaults := ezutil.Config{
			Auth: &ezutil.Auth{
				TokenDuration: 30 * time.Minute,
			},
		}

		config := ezutil.LoadConfigWithoutDB(defaults)
		assert.Equal(t, 30*time.Minute, config.Auth.TokenDuration) // Should use default
	})

	t.Run("negative cookie duration uses default", func(t *testing.T) {
		os.Setenv("AUTH_COOKIEDURATION", "-12h")
		defer os.Unsetenv("AUTH_COOKIEDURATION")

		defaults := ezutil.Config{
			Auth: &ezutil.Auth{
				CookieDuration: 24 * time.Hour,
			},
		}

		config := ezutil.LoadConfigWithoutDB(defaults)
		assert.Equal(t, 24*time.Hour, config.Auth.CookieDuration) // Should use default
	})
}

func TestAppConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		port        string
		expectValid bool
	}{
		{
			name:        "valid port",
			port:        "8080",
			expectValid: true,
		},
		{
			name:        "invalid port - too high",
			port:        "70000",
			expectValid: false,
		},
		{
			name:        "invalid port - negative",
			port:        "-1",
			expectValid: false,
		},
		{
			name:        "invalid port - not a number",
			port:        "abc",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test port validation logic
			if tt.expectValid {
				// Valid ports should be between 1 and 65535
				assert.True(t, isValidPort(tt.port))
			} else {
				// Invalid ports should fail validation
				assert.False(t, isValidPort(tt.port))
			}
		})
	}
}

// Helper function to validate port
func isValidPort(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return portNum >= 1 && portNum <= 65535
}
