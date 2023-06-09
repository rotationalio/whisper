/*
Package config provides settings and configuration for the whipser server by loading the
configuration from the environment and specifying reasonable defaults and required
settings.
*/
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rotationalio/whisper/pkg/logger"
	"github.com/rotationalio/whisper/pkg/sentry"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the whisper service.
type Config struct {
	Maintenance  bool                `split_words:"true" default:"false"`
	Mode         string              `split_words:"true" default:"debug"`
	BindAddr     string              `split_words:"true" required:"false"`
	LogLevel     logger.LevelDecoder `split_words:"true" default:"info"`
	ConsoleLog   bool                `split_words:"true" default:"false"`
	AllowOrigins []string            `split_words:"true" default:"https://whisper.rotational.dev"`
	Google       GoogleConfig
	Sentry       sentry.Config
	processed    bool
}

type GoogleConfig struct {
	Credentials string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS" required:"false"`
	Project     string `envconfig:"GOOGLE_PROJECT_NAME" required:"true"`
	Testing     bool   `split_words:"true" default:"false"`
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("whisper", &conf); err != nil {
		return Config{}, err
	}

	// If the BindAddr is not set, try setting it from $PORT (Google Cloud Run)
	if conf.BindAddr == "" {
		if port := os.Getenv("PORT"); port != "" {
			conf.BindAddr = ":" + port
		}
	}

	// If mode is testing, then google.testing is true, even if it is explicitly set as false.
	if conf.Mode == gin.TestMode {
		conf.Google.Testing = true
	}

	// Extra validation of the configuration
	if err = conf.Validate(); err != nil {
		return Config{}, err
	}

	conf.processed = true
	return conf, nil
}

func (c Config) GetLogLevel() zerolog.Level {
	return zerolog.Level(c.LogLevel)
}

func (c Config) IsZero() bool {
	return !c.processed
}

func (c Config) Validate() error {
	if c.BindAddr == "" {
		return errors.New("must specify either $WHISPER_BIND_ADDR or $PORT")
	}

	if c.Mode != gin.ReleaseMode && c.Mode != gin.DebugMode && c.Mode != gin.TestMode {
		return fmt.Errorf("%q is not a valid gin mode", c.Mode)
	}
	return nil
}
