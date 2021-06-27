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
	"WHISPER_MAINTENANCE":  "false",
	"WHISPER_MODE":         gin.TestMode,
	"WHISPER_BIND_ADDR":    "127.0.0.1:8311",
	"WHISPER_USE_TLS":      "false",
	"WHISPER_DOMAIN":       "localhost",
	"WHISPER_SECRET_KEY":   "supersecretkey",
	"WHISPER_DATABASE_URL": "file::memory:?cache=shared",
	"WHISPER_LOG_LEVEL":    "debug",
	"WHISPER_CONSOLE_LOG":  "false",
}

// WhisperTestSuite mocks the database and gin/http requests for testing endpoints.
type WhisperTestSuite struct {
	suite.Suite
	api    *Server
	conf   config.Config
	router http.Handler
}

func (s *WhisperTestSuite) SetupSuite() {
	// Update the test environment
	for key, val := range testEnv {
		os.Setenv(key, val)
	}

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

func TestWhisper(t *testing.T) {
	suite.Run(t, new(WhisperTestSuite))
}

func (s *WhisperTestSuite) TestGinMode() {
	s.Equal(gin.TestMode, gin.Mode())
}
