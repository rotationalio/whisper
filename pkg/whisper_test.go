package whisper_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/rotationalio/whisper/pkg"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/stretchr/testify/suite"
)

var testEnv = map[string]string{
	"WHISPER_MAINTENANCE":            "false",
	"WHISPER_MODE":                   gin.TestMode,
	"WHISPER_BIND_ADDR":              "127.0.0.1:8311",
	"WHISPER_LOG_LEVEL":              "debug",
	"WHISPER_CONSOLE_LOG":            "false",
	"GOOGLE_APPLICATION_CREDENTIALS": "fixtures/test.json",
	"GOOGLE_PROJECT_NAME":            "test",
}

// WhisperTestSuite mocks the database and gin/http requests for testing endpoints.
type WhisperTestSuite struct {
	suite.Suite
	api     *Server
	conf    config.Config
	router  http.Handler
	prevEnv map[string]string
}

func (s *WhisperTestSuite) SetupSuite() {
	// Store the previous environment to restore test suite
	s.prevEnv = curEnv()

	// Update the test environment
	setEnv()

	// Create test configuration for mocked database and server
	var err error
	s.conf, err = config.New()
	s.NoError(err)

	// Create the api, which will setup both the routes and the database
	s.api, err = New(s.conf)
	s.NoError(err)

	// Get the routes from the server
	s.router = s.api.Routes()

	// Set the server as healthy
	s.api.SetHealth(true)
}

func (s *WhisperTestSuite) TearDownSuite() {
	// Restore the previous environment
	for key, val := range s.prevEnv {
		if val != "" {
			os.Setenv(key, val)
		} else {
			os.Unsetenv(key)
		}
	}
}

func TestWhisper(t *testing.T) {
	suite.Run(t, new(WhisperTestSuite))
}

func (s *WhisperTestSuite) TestGinMode() {
	s.Equal(gin.TestMode, gin.Mode())
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
