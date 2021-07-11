package whisper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
)

var (
	unsuccessful = v1.Reply{Success: false}
	notFound     = v1.Reply{Success: false, Error: "resource not found"}
	notAllowed   = v1.Reply{Success: false, Error: "method not allowed"}
)

// ErrorResponse constructs an new response from the error or returns a success: false.
func ErrorResponse(err interface{}) v1.Reply {
	if err == nil {
		return unsuccessful
	}

	rep := v1.Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}

// NotFound returns a JSON 404 response for the API.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}
