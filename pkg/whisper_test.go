package whisper_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/rotationalio/whisper/pkg"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

func init() {
	// No logging output during tests
	zerolog.SetGlobalLevel(zerolog.PanicLevel)
}

var testEnv = map[string]string{
	"WHISPER_MAINTENANCE":            "false",
	"WHISPER_MODE":                   gin.TestMode,
	"WHISPER_BIND_ADDR":              "127.0.0.1:8311",
	"WHISPER_LOG_LEVEL":              "debug",
	"WHISPER_CONSOLE_LOG":            "false",
	"WHISPER_ALLOW_ORIGINS":          "http://localhost:3000",
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

	// No logging output during tests
	s.conf.LogLevel = config.LogLevelDecoder(zerolog.PanicLevel)

	// Use mock Google Secret Manager for tests
	s.conf.Google.Testing = true

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

func (s *WhisperTestSuite) TestCORS() {
	// We've been having some problems with CORS in the front-end; this test helps
	// ensure that our CORS configuration is correct on the server side.

	// Run an httptest server rather than use the httptest recorder
	server := httptest.NewServer(s.router)
	defer server.Close()

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodOptions, server.URL+"/v1/status", nil)
	s.NoError(err)

	// Add correct origin and headers that we want to test
	req.Header.Add("Origin", "http://localhost:3000")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", "0")
	req.Header.Add("Authorization", "Bearer c2VjcmV0cGFzc3dvcmQ=")

	rep, err := client.Do(req)
	s.NoError(err)

	// The Access Control headers should be present if they're valid
	s.Contains(rep.Header, "Access-Control-Allow-Origin")
	s.Contains(rep.Header, "Access-Control-Allow-Headers")
	s.Contains(rep.Header, "Access-Control-Allow-Methods")

	// The Access-Control-Allow-Origin should match or be * if the CORS policy is valid
	origin := rep.Header.Get("Access-Control-Allow-Origin")
	s.True(origin == "http://localhost:3000" || origin == "*")

	// The Access-Control-Allow-Headers should match our sent headers
	headers := rep.Header.Get("Access-Control-Allow-Headers")
	s.Equal("Origin,Content-Length,Content-Type,Authorization", headers)

	// Add incorrect origin and headers to get CORS rejection
	req, err = http.NewRequest(http.MethodOptions, server.URL+"/v1/status", nil)
	s.NoError(err)

	req.Header.Add("Origin", "http://localhost:666")
	rep, err = client.Do(req)
	s.NoError(err)

	// The Access Control headers should be empty if they're not valid
	s.NotContains(rep.Header, "Access-Control-Allow-Origin")
	s.NotContains(rep.Header, "Access-Control-Allow-Headers")
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
