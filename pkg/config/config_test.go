package config_test

import (
	"os"
	"testing"

	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"WHISPER_MAINTENANCE":  "false",
	"WHISPER_MODE":         "release",
	"WHISPER_BIND_ADDR":    ":443",
	"WHISPER_USE_TLS":      "false",
	"WHISPER_DOMAIN":       "localhost",
	"WHISPER_SECRET_KEY":   "theeaglefliesatmidnight",
	"WHISPER_DATABASE_URL": "postgresql://localhost:5432/whisper",
	"WHISPER_LOG_LEVEL":    "debug",
	"WHISPER_CONSOLE_LOG":  "true",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv()
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})
	setEnv()

	conf, err := config.New()
	require.NoError(t, err)

	// Test configuration set from the environment
	require.Equal(t, false, conf.Maintenance)
	require.Equal(t, testEnv["WHISPER_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.Equal(t, true, conf.ConsoleLog)
}

func TestRequiredConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv("WHISPER_DATABASE_URL", "WHISPER_SECRET_KEY")
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})

	_, err := config.New()
	require.Error(t, err)
	setEnv("WHISPER_DATABASE_URL", "WHISPER_SECRET_KEY")

	conf, err := config.New()
	require.NoError(t, err)

	// Test required configuration
	require.Equal(t, testEnv["WHISPER_DATABASE_URL"], conf.DatabaseURL)
	require.Equal(t, testEnv["WHISPER_SECRET_KEY"], conf.SecretKey)
}

// Returns the current environment for the specified keys, or if no keys are specified
// then returns the current environment for all keys in testEnv.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, envvar := range keys {
			if val, ok := os.LookupEnv(envvar); ok {
				env[envvar] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variable from the testEnv, if no keys are specified, then sets
// all environment variables from the test env.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}
