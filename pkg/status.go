package whisper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
)

const (
	serverStatusOK          = "ok"
	serverStatusNotReady    = "not ready"
	serverStatusUnhealthy   = "unhealthy"
	serverStatusMaintenance = "maintenance"
)

// Status is an unauthenticated endpoint that returns the status of the api server and
// can be used for heartbeats and liveness checks. This status method is the global
// status method, meaning it returns the latest version of the whipser service, no
// matter how many API versions are available.
func (s *Server) Status(c *gin.Context) {
	c.JSON(http.StatusOK, v1.StatusReply{
		Status:  serverStatusOK,
		Uptime:  time.Since(s.started).String(),
		Version: Version(),
	})
}

// Available is middleware that uses the healthy boolean to return a service unavailable
// http status code if the server is shutting down. It does this before all routes to
// ensure that complex handling doesn't bog down the server.
func (s *Server) Available() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check the health and ready status of the server
		s.RLock()
		healthy := s.healthy
		ready := s.ready
		s.RUnlock()

		if s.conf.Maintenance || !healthy || !ready {
			out := v1.StatusReply{
				Uptime:  time.Since(s.started).String(),
				Version: Version(),
			}

			switch {
			case !healthy:
				out.Status = serverStatusUnhealthy
			case !ready:
				out.Status = serverStatusNotReady
			default:
				out.Status = serverStatusMaintenance
			}

			// Write the 503 response and stop processing the request
			c.JSON(http.StatusServiceUnavailable, out)
			c.Abort()
			return
		}

		c.Next()
	}
}

// Healthz is used to alert k8s to the health/liveness status of the server.
func (s *Server) Healthz(c *gin.Context) {
	s.RLock()
	healthy := s.healthy
	s.RUnlock()

	if !healthy {
		c.Data(http.StatusServiceUnavailable, "text/plain", []byte(serverStatusUnhealthy))
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(serverStatusOK))
}

// Readyz is used to alert k8s to the readiness status of the server.
func (s *Server) Readyz(c *gin.Context) {
	s.RLock()
	ready := s.ready
	s.RUnlock()

	if !ready {
		c.Data(http.StatusServiceUnavailable, "text/plain", []byte(serverStatusNotReady))
		return
	}

	c.Data(http.StatusOK, "text/plain", []byte(serverStatusOK))
}
