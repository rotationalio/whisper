package whisper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	serverStatusOK          = "ok"
	serverStatusUnavailable = "unavailable"
)

// Status is an unauthenticated endpoint that returns the status of the api server and
// can be used for heartbeats and liveness checks.
func (s *Server) Status(c *gin.Context) {
	c.JSON(http.StatusOK, StatusResponse{
		Status:    serverStatusOK,
		Timestamp: time.Now(),
		Version:   Version(),
	})
}

// Available is middleware that uses the healthy boolean to return a service unavailable
// http status code if the server is shutting down. It does this before all routes to
// ensure that complex handling doesn't bog down the server.
func (s *Server) Available() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check maintenance mode
		if s.conf.Maintenance {
			c.JSON(http.StatusServiceUnavailable, StatusResponse{
				Status:    "unavailable",
				Error:     "service is currently in maintenance mode",
				Timestamp: time.Now(),
				Version:   Version(),
			})
			c.Abort()
			return
		}

		// Check health status
		s.RLock()
		healthy := s.healthy
		s.RUnlock()

		if !healthy {
			c.JSON(http.StatusServiceUnavailable, StatusResponse{
				Status:    "unavailable",
				Error:     "service is currently shutting down",
				Timestamp: time.Now(),
				Version:   Version(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
