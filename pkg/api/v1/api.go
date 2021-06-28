/*
Package api describes the JSON data structures for v1 requests and responses. It is
stored in its own package to facilitate ease of serialization for Go API clients and to
describe and document the API for external users.
*/
package api

import "time"

//===========================================================================
// Service Interface
//===========================================================================

// Service defines the API, which is implemented by the v1 client.
type Service interface {
	Status() (*StatusReply, error)
}

//===========================================================================
// Top Level Requests and Responses
//===========================================================================

// Reply contains standard fields that are embedded in most API responses
type Reply struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

// StatusReply is returned on status requests. Note that no request is needed.
type StatusReply struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Version   string    `json:"version,omitempty"`
	Error     string    `json:"error,omitempty" yaml:"error,omitempty"`
}
