package whisper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	unsuccessful = Response{Success: false}
	notFound     = Response{Success: false, Error: "resource not found"}
	notAllowed   = Response{Success: false, Error: "method not allowed"}
)

// ErrorResponse constructs an new response from the error or returns a success: false.
func ErrorResponse(err error) Response {
	if err == nil {
		return unsuccessful
	}
	return Response{Success: false, Error: err.Error()}
}

// NotFound returns a JSON 404 response for the API.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}
