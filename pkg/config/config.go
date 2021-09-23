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
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the whisper service.
type Config struct {
	Maintenance  bool            `split_words:"true" default:"false"`
	Mode         string          `split_words:"true" default:"debug"`
	BindAddr     string          `split_words:"true" required:"false"`
	LogLevel     LogLevelDecoder `split_words:"true" default:"info"`
	ConsoleLog   bool            `split_words:"true" default:"false"`
	AllowOrigins []string        `split_words:"true" default:"https://whisper.rotational.dev"`
	Google       GoogleConfig
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

// LogLevelDecoder deserializes the log level from a config string.
type LogLevelDecoder zerolog.Level

// Decode implements envconfig.Decoder
func (ll *LogLevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "panic":
		*ll = LogLevelDecoder(zerolog.PanicLevel)
	case "fatal":
		*ll = LogLevelDecoder(zerolog.FatalLevel)
	case "error":
		*ll = LogLevelDecoder(zerolog.ErrorLevel)
	case "warn":
		*ll = LogLevelDecoder(zerolog.WarnLevel)
	case "info":
		*ll = LogLevelDecoder(zerolog.InfoLevel)
	case "debug":
		*ll = LogLevelDecoder(zerolog.DebugLevel)
	case "trace":
		*ll = LogLevelDecoder(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log level %q", value)
	}
	return nil
}
