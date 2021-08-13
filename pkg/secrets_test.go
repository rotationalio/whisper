package whisper_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/rotationalio/whisper/pkg"
	"github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/stretchr/testify/require"
)

func TestParseBearerToken(t *testing.T) {
	password := base64.URLEncoding.EncodeToString([]byte("supersecretsquirrel"))
	tt := []struct {
		header   string
		expected string
	}{
		// Success cases
		{"Bearer " + password, "supersecretsquirrel"},
		{"bearer " + password, "supersecretsquirrel"},
		{"   Bearer    " + password, "supersecretsquirrel"},

		// Failure cases
		{password, ""},                        // No bearer token
		{"Bearer supersecretsquirrel", ""},    // Not base64 encoded
		{"weird foo string with nothing", ""}, // No bearer realm
	}

	for _, tc := range tt {
		require.Equal(t, tc.expected, ParseBearerToken(tc.header))
	}
}

func (s *WhisperTestSuite) TestCreateFetchSecret() {
	// Send a create request and get create reply
	rep1 := s.sendCreateSecret(&api.CreateSecretRequest{
		Secret:   "do not share this with anyone",
		Password: "",
		Accesses: 1,
		Lifetime: api.Duration(30 * time.Minute),
		Filename: "",
		IsBase64: false,
	}, http.StatusCreated)
	s.NotEmpty(rep1)

	// Send a fetch request and get a fetch reply
	rep2 := s.sendFetchRequest(rep1.Token, "", http.StatusOK)
	s.NotEmpty(rep2)
	s.Equal("do not share this with anyone", rep2.Secret)
	s.Equal("", rep2.Filename)
	s.False(rep2.IsBase64)
	s.NotZero(rep2.Created)
	s.Equal(1, rep2.Accesses)
	s.True(rep2.Destroyed)

	// Next request should return 404 since the secret should have been destroyed
	s.sendFetchRequest(rep1.Token, "", http.StatusNotFound)
}

// TODO: CreateFetchSecretPasswordFlow
// TODO: CreateDeleteSecretFlow
// TODO: CreateDeleteSecretPassword Flow

func (s *WhisperTestSuite) sendCreateSecret(in *api.CreateSecretRequest, code int) (out *api.CreateSecretReply) {
	indata, err := json.Marshal(in)
	s.NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/v1/secrets", bytes.NewReader(indata))
	req.Header.Add("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	rep := w.Result()
	defer rep.Body.Close()

	s.Equal(code, rep.StatusCode)

	out = &api.CreateSecretReply{}
	s.NoError(json.NewDecoder(rep.Body).Decode(&out))
	return out
}

func (s *WhisperTestSuite) sendFetchRequest(token, password string, code int) *api.FetchSecretReply {
	path := fmt.Sprintf("/v1/secrets/%s", token)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)

	// TODO: add password header

	s.router.ServeHTTP(w, req)

	rep := w.Result()
	defer rep.Body.Close()

	s.Equal(code, rep.StatusCode)

	out := &api.FetchSecretReply{}
	s.NoError(json.NewDecoder(rep.Body).Decode(&out))
	return out
}
