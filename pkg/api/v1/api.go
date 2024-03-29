/*
Package api describes the JSON data structures for v1 requests and responses. It is
stored in its own package to facilitate ease of serialization for Go API clients and to
describe and document the API for external users.
*/
package api

import (
	"context"
	"time"
)

//===========================================================================
// Service Interface
//===========================================================================

// Service defines the API, which is implemented by the v1 client.
type Service interface {
	Status(ctx context.Context) (out *StatusReply, err error)
	CreateSecret(ctx context.Context, in *CreateSecretRequest) (out *CreateSecretReply, err error)
	FetchSecret(ctx context.Context, token, password string) (out *FetchSecretReply, err error)
	DestroySecret(ctx context.Context, token, password string) (out *DestroySecretReply, err error)
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
	Status  string `json:"status"`
	Uptime  string `json:"uptime,omitempty"`
	Version string `json:"version,omitempty"`
	Error   string `json:"error,omitempty" yaml:"error,omitempty"`
}

//===========================================================================
// Secret REST API
//===========================================================================

type CreateSecretRequest struct {
	Secret   string   `json:"secret" binding:"required"` // the secret can be a string of any length or base64 encoded data
	Password string   `json:"password,omitempty"`        // a password that must be used to retrieve the secret
	Accesses int      `json:"accesses,omitempty"`        // specify the number of times the secret can be accessed; default is 1, if negative, can be accessed until the secret expires
	Lifetime Duration `json:"lifetime,omitempty"`        // how long the secret will last before being deleted
	Filename string   `json:"filename,omitempty"`        // if the secret is a filename, the name of the file
	IsBase64 bool     `json:"is_base64"`                 // if the secret is base64 encoded or not
}

type CreateSecretReply struct {
	Token   string    `json:"token"`   // the token used to retrieve the secret (so the URL doesn't have to be parsed)
	Expires time.Time `json:"expires"` // the timestamp when the secret will have expired
}

type FetchSecretReply struct {
	Secret    string    `json:"secret"`             // the secret retrieved by the database, which is now deleted
	Filename  string    `json:"filename,omitempty"` // the name of the file used to create the secret to save as a file
	IsBase64  bool      `json:"is_base64"`          // if the secret is base64 encoded data
	Created   time.Time `json:"created"`            // the timestamp the secret was created
	Accesses  int       `json:"accesses"`           // the number of times the secret has been accessed
	Destroyed bool      `json:"destroyed"`          // if the secret was destroyed after the fetch
}

type DestroySecretReply struct {
	Destroyed bool `json:"destroyed"` // if the secret was destroyed or not
}
