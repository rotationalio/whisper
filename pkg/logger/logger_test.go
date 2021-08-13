package logger_test

import (
	"os"

	. "github.com/rotationalio/whisper/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ExampleLogger() {
	// Initialize zerolog with GCP logging requirements
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = "foo" // This should be RFC3339 but blanked for testing purposes
	zerolog.TimestampFieldName = GCPFieldKeyTime
	zerolog.MessageFieldName = GCPFieldKeyMsg
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Add the severity hook for GCP logging
	var gcpHook SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()

	log.Trace().Msg("super verbose")
	log.Debug().Msg("nothing to see here")
	log.Info().Msg("hello world")
	log.Warn().Msg("be careful")
	log.Error().Msg("something bad happend")
	// Output:
	// {"level":"info","severity":"INFO","time":"foo","message":"hello world"}
	// {"level":"warn","severity":"WARNING","time":"foo","message":"be careful"}
	// {"level":"error","severity":"ERROR","time":"foo","message":"something bad happend"}
}
