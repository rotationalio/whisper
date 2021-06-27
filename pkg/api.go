package whisper

import "time"

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Response contains standard fields that are embedded in most API responses
type Response struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusResponse is returned on status requests. Note that no request is needed.
type StatusResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Version   string    `json:"version,omitempty"`
	Error     string    `json:"error,omitempty" yaml:"error,omitempty"`
}
