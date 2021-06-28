/*
Package config provides settings and configuration for the whipser server by loading the
configuration from the environment and specifying reasonable defaults and required
settings.
*/
package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

// Config uses envconfig to load required settings from the environment and validate
// them in preparation for running the whisper service.
type Config struct {
	Maintenance bool            `split_words:"true" default:"false"`
	Mode        string          `split_words:"true" default:"debug"`
	BindAddr    string          `split_words:"true" default:":8318"`
	UseTLS      bool            `split_words:"true" default:"false"`
	Domain      string          `split_words:"true" default:"localhost"`
	SecretKey   string          `split_words:"true" required:"true"`
	DatabaseURL string          `split_words:"true" required:"true"`
	LogLevel    LogLevelDecoder `split_words:"true" default:"info"`
	ConsoleLog  bool            `split_words:"true" default:"false"`
	processed   bool
}

// New creates a new Config object, loading environment variables and defaults.
func New() (_ Config, err error) {
	var conf Config
	if err = envconfig.Process("whisper", &conf); err != nil {
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
