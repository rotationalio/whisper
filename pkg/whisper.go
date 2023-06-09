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

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rotationalio/whisper/pkg/logger"
	"github.com/rotationalio/whisper/pkg/sentry"
	"github.com/rotationalio/whisper/pkg/vault"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// Initialize zerolog with GCP logging requirements
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFieldName = logger.GCPFieldKeyTime
	zerolog.MessageFieldName = logger.GCPFieldKeyMsg

	// Add the severity hook for GCP logging
	var gcpHook logger.SeverityHook
	log.Logger = zerolog.New(os.Stdout).Hook(gcpHook).With().Timestamp().Logger()
}

const ServiceName = "whisper"

func New(conf config.Config) (s *Server, err error) {
	// Load the default configuration from the environment
	if conf.IsZero() {
		if conf, err = config.New(); err != nil {
			return nil, err
		}
	}

	// Set the global level
	zerolog.SetGlobalLevel(zerolog.Level(conf.LogLevel))

	// Configure Sentry
	if conf.Sentry.UseSentry() {
		// Set the release version if not already set
		if conf.Sentry.Release == "" {
			conf.Sentry.Release = Version()
		}

		if err = sentry.Init(conf.Sentry); err != nil {
			return nil, err
		}
	}

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
	// Note that if conf.Google.Testing is true, a mock secret manager will be created
	if s.vault, err = vault.New(conf.Google); err != nil {
		return nil, err
	}
	log.Debug().Msg("connected to google secret manager")

	// Create the Gin router and setup its routes
	gin.SetMode(conf.Mode)
	s.router = gin.New()
	s.router.RedirectTrailingSlash = true
	s.router.RedirectFixedPath = false
	s.router.HandleMethodNotAllowed = true
	s.router.ForwardedByClientIP = true
	s.router.UseRawPath = false
	s.router.UnescapePathValues = true
	if err = s.setupRoutes(); err != nil {
		return nil, err
	}

	// Create the http server
	s.srv = &http.Server{
		Addr:         s.conf.BindAddr,
		Handler:      s.router,
		ErrorLog:     nil,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Debug().Msg("created http server with gin router")
	return s, nil
}

type Server struct {
	sync.RWMutex
	conf    config.Config        // configuration of the API server
	srv     *http.Server         // handle to a custom http server with specified API defaults
	router  *gin.Engine          // the http handler and associated middlware
	vault   *vault.SecretManager // storage for all secrets the whisper application manages
	healthy bool                 // application state of the server for health checks
	ready   bool                 // application state of the server for ready checks
	started time.Time            // the timestamp when the server was started
	errc    chan error           // synchronize shutdown gracefully
}

func (s *Server) Serve() (err error) {
	s.osSignals()
	s.SetStatus(true, true)

	if s.conf.Maintenance {
		log.Warn().Msg("starting server in maintenance mode")
	}

	s.started = time.Now()
	log.Info().Str("addr", s.conf.BindAddr).Msg("whisper server started")

	if err = s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	if err = <-s.errc; err != nil {
		sentry.Fatal(nil).Err(err).Msg("fatal error, server stopped")
	}
	return nil
}

func (s *Server) Shutdown() (err error) {
	log.Info().Msg("gracefully shutting down whisper server")
	s.SetStatus(false, false)

	errs := make([]error, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	s.srv.SetKeepAlivesEnabled(false)
	if err = s.srv.Shutdown(ctx); err != nil {
		sentry.Error(nil).Err(err).Msg("could not shutdown http server")
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
	// Instantiate Sentry handlers
	var tags gin.HandlerFunc
	if s.conf.Sentry.UseSentry() {
		tagmap := map[string]string{"service": ServiceName}
		tags = sentry.UseTags(tagmap)
	}

	var tracing gin.HandlerFunc
	if s.conf.Sentry.UsePerformanceTracking() {
		tagmap := map[string]string{"service": ServiceName}
		tracing = sentry.TrackPerformance(tagmap)
	}

	// Setup CORS configuration
	corsConf := cors.Config{
		AllowOrigins:     s.conf.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-TOKEN", "sentry-trace", "baggage"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Application Middleware
	// NOTE: ordering is important to how middleware is handled
	middlewares := []gin.HandlerFunc{
		// Logging should be on the outside so we can record the correct latency of requests
		// NOTE: logging panics will not recover
		logger.GinLogger(ServiceName, Version()),

		// Panic recovery middleware
		gin.Recovery(),
		sentrygin.New(sentrygin.Options{
			Repanic:         true,
			WaitForDelivery: false,
		}),

		// Add searchable tags to sentry context
		tags,

		// Tracing helps us measure performance metrics with Sentry
		tracing,

		// CORS configuration allows the front-end to make cross-origin requests
		cors.New(corsConf),

		// Mainenance mode handling
		s.Available(),
	}

	// Add the middleware to the router
	for _, middleware := range middlewares {
		if middleware != nil {
			s.router.Use(middleware)
		}
	}

	// Redirect the root to the current version root
	s.router.GET("/", s.RedirectVersion)

	// Add the v1 API routes (currently the only version)
	v1 := s.router.Group("/v1")
	{
		// Heartbeat route
		v1.GET("/status", s.Status)

		// Secrets REST resource
		v1.POST("/secrets", s.CreateSecret)
		v1.GET("/secrets/:token", s.FetchSecret)
		v1.DELETE("/secrets/:token", s.DestroySecret)
	}

	// Kubernetes liveness probes
	s.router.GET("/healthz", s.Healthz)
	s.router.GET("/livez", s.Healthz)
	s.router.GET("/readyz", s.Readyz)

	// NotFound and NotAllowed requests
	s.router.NoRoute(NotFound)
	s.router.NoMethod(NotAllowed)
	return nil
}

// SetHealth sets the health status on the API server, putting it into unavailable mode
// if health is false, and removing maintenance mode if health is true. Here primarily
// for testing purposes since it is unlikely an outside caller can access this.
func (s *Server) SetStatus(health, ready bool) {
	s.Lock()
	s.healthy = health
	s.ready = ready
	s.Unlock()
	log.Debug().Bool("health", health).Bool("ready", ready).Msg("server status set")
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
