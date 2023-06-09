package sentry

import (
	"errors"

	"github.com/getsentry/sentry-go"
)

// Sentry configuration for use in application-configuration
type Config struct {
	DSN              string  `split_words:"true"`
	ServerName       string  `split_words:"true"`
	Environment      string  `split_words:"true" default:"production"`
	Release          string  `split_words:"true"`
	TrackPerformance bool    `split_words:"true" default:"true"`
	SampleRate       float64 `split_words:"true" default:"0.25"`
	ReportErrors     bool    `split_words:"true" default:"true"`
	Repanic          bool    `ignored:"true"`
	Debug            bool    `default:"false"`
}

// Returns true if Sentry is enabled (e.g. a DSN is configured)
func (c Config) UseSentry() bool {
	return c.DSN != ""
}

// Performance tracking is enabled if Sentry is enabled and track performance is explicitly set
func (c Config) UsePerformanceTracking() bool {
	return c.UseSentry() && c.TrackPerformance
}

func (c Config) Validate() error {
	if c.UseSentry() && c.Environment == "" {
		return errors.New("invalid configuration: environment must be configured when Sentry is enabled")
	}
	return nil
}

func (c Config) ClientOptions() sentry.ClientOptions {
	return sentry.ClientOptions{
		Dsn:              c.DSN,
		Environment:      c.Environment,
		Release:          c.Release,
		AttachStacktrace: true,
		Debug:            c.Debug,
		ServerName:       c.ServerName,
		EnableTracing:    c.TrackPerformance,
		TracesSampleRate: c.SampleRate,
	}
}
