package logger

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

// LogLevelDecoder deserializes the log level from a config string.
type LevelDecoder zerolog.Level

// Names of log levels for use in encoding/decoding from strings.
const (
	llPanic = "panic"
	llFatal = "fatal"
	llError = "error"
	llWarn  = "warn"
	llInfo  = "info"
	llDebug = "debug"
	llTrace = "trace"
)

// Decode implements envconfig.Decoder
func (ll *LevelDecoder) Decode(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case llPanic:
		*ll = LevelDecoder(zerolog.PanicLevel)
	case llFatal:
		*ll = LevelDecoder(zerolog.FatalLevel)
	case llError:
		*ll = LevelDecoder(zerolog.ErrorLevel)
	case llWarn:
		*ll = LevelDecoder(zerolog.WarnLevel)
	case llInfo:
		*ll = LevelDecoder(zerolog.InfoLevel)
	case llDebug:
		*ll = LevelDecoder(zerolog.DebugLevel)
	case llTrace:
		*ll = LevelDecoder(zerolog.TraceLevel)
	default:
		return fmt.Errorf("unknown log level %q", value)
	}
	return nil
}

// Encode converts the loglevel into a string for use in YAML and JSON
func (ll *LevelDecoder) Encode() (string, error) {
	switch zerolog.Level(*ll) {
	case zerolog.PanicLevel:
		return llPanic, nil
	case zerolog.FatalLevel:
		return llFatal, nil
	case zerolog.ErrorLevel:
		return llError, nil
	case zerolog.WarnLevel:
		return llWarn, nil
	case zerolog.InfoLevel:
		return llInfo, nil
	case zerolog.DebugLevel:
		return llDebug, nil
	case zerolog.TraceLevel:
		return llTrace, nil
	default:
		return "", fmt.Errorf("unknown log level %d", ll)
	}
}

func (ll LevelDecoder) String() string {
	ls, _ := ll.Encode()
	return ls
}
