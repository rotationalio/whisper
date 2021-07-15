package whisper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rotationalio/whisper/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Initialize zerolog with GCP logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(zerolog.Level(conf.LogLevel))

	// Set human readable logging if specified
	if conf.ConsoleLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Check that a cryptographically secure PNRG is available
	if err = checkAvailablePRNG(); err != nil {
		return nil, err
	}

	// Create the server and prepare to serve
	s = &Server{conf: conf, errc: make(chan error, 1), healthy: false}

	// Create the vault to store secrets in (Google Secret Manager)
	// TODO: if we're in test mode, create a mock secret manager
	if conf.Mode != gin.TestMode {
		if s.vault, err = NewSecretManager(conf.Google); err != nil {
			return nil, err
		}
		log.Debug().Msg("connected to google secret manager")
	}

	// Create the router
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	s.router.Use(ginzerolog.Logger("gin"))
	s.router.Use(gin.Recovery())
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:         s.conf.BindAddr,
		Handler:      s.router,
		ErrorLog:     nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Debug().Msg("created http server with gin router")
	return s, nil
}

type Server struct {
	sync.RWMutex
	conf    config.Config  // configuration of the API server
	srv     *http.Server   // handle to a custom http server with specified API defaults
	router  *gin.Engine    // the http handler and associated middlware
	vault   *SecretManager // storage for all secrets the whisper application manages
	healthy bool           // application state of the server
	errc    chan error     // synchronize shutdown gracefully
}

func (s *Server) Serve() (err error) {
	s.SetHealth(!s.conf.Maintenance)
	s.osSignals()

	if s.conf.Maintenance {
		log.Warn().Msg("starting server in maintenance mode")
	}

	log.Info().Str("addr", s.conf.BindAddr).Msg("whisper server started")
	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	if err = <-s.errc; err != nil {
		log.Error().Err(err).Msg("fatal error, server stopped")
	}
	return nil
}

func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down whisper server")
	s.SetHealth(false)

	errs := make([]error, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("could not shutdown http server")
		errs = append(errs, err)
	}

	switch len(errs) {
	case 0:
		log.Debug().Msg("successful shutdown of whisper server")
		close(s.errc)
		return nil
	case 1:
		s.errc <- errors.New("shutdown failed with error")
		return errs[0]
	default:
		s.errc <- errors.New("shutdown failed with multiple errors")
		return fmt.Errorf("%d errors occurred during shutdown", len(errs))
	}
}

// Routes returns the API router and is primarily exposed for testing purposes.
func (s *Server) Routes() http.Handler {
	return s.router
}

func (s *Server) setupRoutes() (err error) {

	// Application Middleware
	s.router.Use(s.Available())

	// Redirect the root to the current version root
	s.router.GET("/", s.RedirectVersion)

	// Add the v1 API routes (currently the only version)
	v1 := s.router.Group("/v1")
	v1.GET("/status", s.Status)
	v1.POST("/secrets", s.CreateSecret)
	v1.GET("/secrets/:token", s.FetchSecret)
	v1.DELETE("/secrets/:token", s.DestroySecret)

	// NotFound and NotAllowed requests
	s.router.NoRoute(NotFound)
	s.router.NoMethod(NotAllowed)
	return nil
}

// SetHealth sets the health status on the API server, putting it into unavailable mode
// if health is false, and removing maintenance mode if health is true. Here primarily
// for testing purposes since it is unlikely an outside caller can access this.
func (s *Server) SetHealth(health bool) {
	s.Lock()
	s.healthy = health
	s.Unlock()
	log.Debug().Bool("health", health).Msg("server health set")
}

func (s *Server) osSignals() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		s.Shutdown()
	}()
	log.Debug().Msg("listening for OS signals SIGINT and SIGTERM")
}
