package whisper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
)

const (
	serverStatusOK          = "ok"
	serverStatusUnavailable = "unavailable"
)

// Status is an unauthenticated endpoint that returns the status of the api server and
// can be used for heartbeats and liveness checks. This status method is the global
// status method, meaning it returns the latest version of the whipser service, no
// matter how many API versions are available.
func (s *Server) Status(c *gin.Context) {
	c.JSON(http.StatusOK, v1.StatusReply{
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
		// Check health status (if unhealthy, assume maintenance mode)
		s.RLock()
		if !s.healthy {
			c.JSON(http.StatusServiceUnavailable, v1.StatusReply{
				Status:    serverStatusUnavailable,
				Error:     "service is currently in maintenance mode",
				Timestamp: time.Now(),
				Version:   Version(),
			})
			c.Abort()
			s.RUnlock()
			return
		}
		s.RUnlock()
		c.Next()
	}
}
