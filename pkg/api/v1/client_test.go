package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			require.Equal(t, int64(0), r.ContentLength)
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "{\"hello\":\"world\"}")
			return
		}

		require.Equal(t, int64(18), r.ContentLength)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "{\"error\":\"bad request\"}")
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	// Ensure that the latest version of the client is returned
	apiv1, ok := client.(*api.APIv1)
	require.True(t, ok)

	// Create a new GET request to a basic path
	req, err := apiv1.NewRequest(context.TODO(), http.MethodGet, "/foo", nil)
	require.NoError(t, err)

	require.Equal(t, "/foo", req.URL.Path)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "Whisper/1.0", req.Header.Get("User-Agent"))
	require.Equal(t, "application/json", req.Header.Get("Accept"))
	require.Equal(t, "application/json", req.Header.Get("Content-Type"))

	data := make(map[string]string)
	rep, err := apiv1.Do(req, &data, true)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rep.StatusCode)
	require.Contains(t, data, "hello")
	require.Equal(t, "world", data["hello"])

	// Create a new POST request and check error handling
	req, err = apiv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data)
	require.NoError(t, err)
	rep, err = apiv1.Do(req, nil, false)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rep.StatusCode)

	req, err = apiv1.NewRequest(context.TODO(), http.MethodPost, "/bar", data)
	require.NoError(t, err)
	_, err = apiv1.Do(req, nil, true)
	require.EqualError(t, err, "[400] 400 Bad Request")
}

func TestStatus(t *testing.T) {
	fixture := &api.StatusReply{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.test",
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/status", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.Status(context.TODO())
	require.NoError(t, err)
	require.Equal(t, fixture.Status, out.Status)
	require.True(t, fixture.Timestamp.Equal(out.Timestamp))
	require.Equal(t, fixture.Version, out.Version)
	require.Empty(t, out.Error)
}

func TestCreateSecret(t *testing.T) {
	fixture := &api.CreateSecretReply{
		Token:   "abc1234cde",
		Expires: time.Now().Add(24 * time.Hour),
	}

	req := &api.CreateSecretRequest{
		Secret:   "super secret squirrel",
		Password: "unlockingkey",
		Accesses: 1,
		Lifetime: api.Duration(time.Hour * 24),
		Filename: "",
		IsBase64: false,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/secrets", r.URL.Path)

		// Must be able to deserialize the request
		in := new(api.CreateSecretRequest)
		err := json.NewDecoder(r.Body).Decode(in)
		require.NoError(t, err)

		require.Equal(t, req.Secret, in.Secret)
		require.Equal(t, req.Password, in.Password)
		require.Equal(t, req.Accesses, in.Accesses)
		require.Equal(t, req.Lifetime, in.Lifetime)
		require.Equal(t, req.Filename, in.Filename)
		require.Equal(t, req.IsBase64, in.IsBase64)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.CreateSecret(context.TODO(), req)
	require.NoError(t, err)
	require.Equal(t, fixture.Token, out.Token)
	require.True(t, fixture.Expires.Equal(out.Expires))
}

func TestFetchSecretNoPassword(t *testing.T) {
	fixture := &api.FetchSecretReply{
		Secret:    "the eagle flies at midnight",
		Filename:  "",
		IsBase64:  false,
		Created:   time.Now().Add(-3 * time.Hour),
		Accesses:  1,
		Destroyed: true,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/secrets/abcd1234dcba", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.FetchSecret(context.TODO(), "abcd1234dcba", "")
	require.NoError(t, err)
	require.Equal(t, fixture.Secret, out.Secret)
	require.Equal(t, fixture.Filename, out.Filename)
	require.Equal(t, fixture.IsBase64, out.IsBase64)
	require.True(t, fixture.Created.Equal(out.Created))
	require.Equal(t, fixture.Accesses, out.Accesses)
	require.Equal(t, fixture.Destroyed, out.Destroyed)
}

func TestFetchSecretPassword(t *testing.T) {
	fixture := &api.FetchSecretReply{}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/secrets/abcd1234dcba", r.URL.Path)
		require.Equal(t, "Bearer c3VwZXJzZWNyZXQ=", r.Header.Get("Authorization"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	_, err = client.FetchSecret(context.TODO(), "abcd1234dcba", "supersecret")
	require.NoError(t, err)
}

func TestDestroySecretNoPassword(t *testing.T) {
	fixture := &api.DestroySecretReply{
		Destroyed: true,
	}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/secrets/abcd1234dcba", r.URL.Path)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	out, err := client.DestroySecret(context.TODO(), "abcd1234dcba", "")
	require.NoError(t, err)
	require.Equal(t, fixture.Destroyed, out.Destroyed)
}

func TestDestroySecretPassword(t *testing.T) {
	fixture := &api.DestroySecretReply{}

	// Create a Test Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/secrets/abcd1234dcba", r.URL.Path)
		require.Equal(t, "Bearer c3VwZXJzZWNyZXQ=", r.Header.Get("Authorization"))

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fixture)
	}))
	defer ts.Close()

	// Create a Client that makes requests to the test server
	client, err := api.New(ts.URL)
	require.NoError(t, err)

	_, err = client.DestroySecret(context.TODO(), "abcd1234dcba", "supersecret")
	require.NoError(t, err)
}
