package config_test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

var testEnv = map[string]string{
	"WHISPER_MAINTENANCE":            "false",
	"WHISPER_MODE":                   "release",
	"WHISPER_BIND_ADDR":              ":443",
	"WHISPER_LOG_LEVEL":              "debug",
	"WHISPER_CONSOLE_LOG":            "true",
	"GOOGLE_APPLICATION_CREDENTIALS": "fixtures/whisper-sa.json",
	"GOOGLE_PROJECT_NAME":            "test-project",
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
	require.Equal(t, gin.ReleaseMode, conf.Mode)
	require.Equal(t, testEnv["WHISPER_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, zerolog.DebugLevel, conf.GetLogLevel())
	require.Equal(t, testEnv["GOOGLE_APPLICATION_CREDENTIALS"], conf.Google.Credentials)
	require.Equal(t, testEnv["GOOGLE_PROJECT_NAME"], conf.Google.Project)
	require.Equal(t, true, conf.ConsoleLog)
}

func TestRequiredConfig(t *testing.T) {
	// Set required environment variables and cleanup after
	prevEnv := curEnv("WHISPER_BIND_ADDR", "GOOGLE_PROJECT_NAME")
	t.Cleanup(func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	})

	// Required EnvVars from struct tags
	conf, err := config.New()
	require.Error(t, err)
	require.True(t, conf.IsZero())
	setEnv("GOOGLE_PROJECT_NAME")

	// Required EnvVars from Validate
	conf, err = config.New()
	require.Error(t, err)
	require.True(t, conf.IsZero())
	setEnv("WHISPER_BIND_ADDR")

	conf, err = config.New()
	require.NoError(t, err)
	require.False(t, conf.IsZero())

	// Test required configuration
	require.Equal(t, testEnv["WHISPER_BIND_ADDR"], conf.BindAddr)
	require.Equal(t, testEnv["GOOGLE_PROJECT_NAME"], conf.Google.Project)

	// Test the use of $PORT instead of WHISPER_BIND_ADDR
	os.Unsetenv("WHISPER_BIND_ADDR")
	os.Setenv("PORT", "5356")
	conf, err = config.New()
	require.NoError(t, err)
	require.False(t, conf.IsZero())
	require.Equal(t, ":5356", conf.BindAddr)
}

func TestLogLevelDecoder(t *testing.T) {
	tt := []struct {
		name  string
		level zerolog.Level
	}{
		{"panic", zerolog.PanicLevel},
		{"FATAL", zerolog.FatalLevel},
		{"  eRrOr  ", zerolog.ErrorLevel},
		{"warn", zerolog.WarnLevel},
		{"info", zerolog.InfoLevel},
		{"DEBUG", zerolog.DebugLevel},
		{"  trace", zerolog.TraceLevel},
	}

	for _, tc := range tt {
		ll := new(config.LogLevelDecoder)
		require.NoError(t, ll.Decode(tc.name))
		require.Equal(t, tc.level, zerolog.Level(*ll))
	}

	// Handle unknown level
	ll := new(config.LogLevelDecoder)
	require.EqualError(t, ll.Decode("foo"), "unknown log level \"foo\"")
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
