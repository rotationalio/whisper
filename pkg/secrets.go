package whisper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/rs/zerolog/log"
)

// DefaultSecretLifetime is one week after which the secret will be deleted.
const DefaultSecretLifetime = time.Hour * 24 * 7

var tmpSecretsStore = make(map[string]v1.CreateSecretRequest)

// CreateSecret handles an incoming CreateSecretRequest and attempts to create a new
// secret that will only be displayed when the correct link is retrieved.
func (s *Server) CreateSecret(c *gin.Context) {
	// Parse incoming JSON data from the client request
	var req v1.CreateSecretRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("invalid create secret request"))
		return
	}

	// Create the reply to return to the User
	var err error
	rep := &v1.CreateSecretReply{}

	// Make a random URL to store the secret in
	if rep.Token, err = GenerateUniqueURL(); err != nil {
		log.Error().Err(err).Msg("could not generate unique token for secret")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// Compute the expiration time from the request
	if req.Lifetime == v1.Duration(0) {
		fmt.Println(req.Lifetime)
		rep.Expires = time.Now().Add(DefaultSecretLifetime)
	} else {
		rep.Expires = time.Now().Add(time.Duration(req.Lifetime))
	}

	// Store the password in the database
	tmpSecretsStore[rep.Token] = req

	// Return success
	c.JSON(http.StatusCreated, rep)
}

// FetchSecret handles an incoming FetchSecretRequest and attempts to retrieve the
// secret from the database and return it to the user. This function also handles the
// password and ensures that a 404 is returned to obfuscate the existence of the secret
// on bad requests.
// TODO: handle passwords (should this also be post?)
// TODO: add created timestamp
func (s *Server) FetchSecret(c *gin.Context) {
	token := c.Param("token")
	req, ok := tmpSecretsStore[token]
	if !ok {
		c.JSON(http.StatusNotFound, ErrorResponse("secret not found"))
		return
	}

	rep := v1.FetchSecretReply{
		Secret:   req.Secret,
		Filename: req.Filename,
		IsBase64: req.IsBase64,
	}
	c.JSON(http.StatusOK, rep)
}

const (
	generateUniqueLength   = 32
	generateUniqueAttempts = 8
)

// GenerateUniqueURL is a helper function that uses crypto/rand to create a random
// URL-safe string and determines if it is in the database or not. If it finds a
// collision it attempts to find a unique string for a fixed number of attempts before
// quitting.
func GenerateUniqueURL() (token string, err error) {
	for i := 0; i < generateUniqueAttempts; i++ {
		// Create a random array of bytes
		buf := make([]byte, generateUniqueLength)
		if _, err = rand.Read(buf); err != nil {
			return "", err
		}

		// Convert random array into URL safe base64 encoded string
		token := base64.RawURLEncoding.EncodeToString(buf)

		// Check if the token exists already in the database
		if _, ok := tmpSecretsStore[token]; !ok {
			return token, nil
		}
	}
	return "", fmt.Errorf("could not generate unique URL after %d attempts", generateUniqueAttempts)
}

// Check that a cryptographically secure PRNG is available.
func checkAvailablePRNG() (err error) {
	buf := make([]byte, 1)
	if _, err = io.ReadFull(rand.Reader, buf); err != nil {
		return fmt.Errorf("crypto/rand is unavailable: failed with %#v", err)
	}
	return nil
}
