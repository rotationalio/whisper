package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func New(endpoint string) (_ Service, err error) {
	c := &api{
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Second,
		},
	}
	if c.endpoint, err = url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %s", err)
	}
	return c, nil
}

// API implements the Service interface.
// TODO: add redirect check that ensures the client only accesses v1 routes.
type api struct {
	endpoint *url.URL
	client   *http.Client
}

func (s api) Status() (_ *StatusReply, err error) {
	// Resolve the URL reference
	u := s.endpoint.ResolveReference(&url.URL{Path: "/v1/status"})

	//  Make the HTTP request
	var rep *http.Response
	if rep, err = s.client.Get(u.String()); err != nil {
		return nil, fmt.Errorf("could not execute request: %s", err)
	}

	// Detect other errors
	if rep.StatusCode != http.StatusOK && rep.StatusCode != http.StatusServiceUnavailable {
		return nil, fmt.Errorf("[%d] %s", rep.StatusCode, rep.Status)
	}

	// Deserialize the JSON data from the response
	status := &StatusReply{}
	if err = json.NewDecoder(rep.Body).Decode(status); err != nil {
		return nil, fmt.Errorf("could not deserialize StatusReply: %s", err)
	}
	return status, nil
}
