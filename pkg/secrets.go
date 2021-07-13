package whisper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	v1 "github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/rs/zerolog/log"
)

// DefaultSecretLifetime is one week after which the secret will be destroyed.
const DefaultSecretLifetime = time.Hour * 24 * 7

// DefaultSecretAccesses ensures that once the secret is fetched it is destroyed
const DefaultSecretAccesses = 1

var tmpSecretsStore = make(map[string]string)
var tmpSecretsMeta = make(map[string]*SecretMetadata)

// CreateSecret handles an incoming CreateSecretRequest and attempts to create a new
// secret that will only be displayed when the correct link is retrieved.
func (s *Server) CreateSecret(c *gin.Context) {
	// Parse incoming JSON data from the client request
	var req v1.CreateSecretRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse("invalid create secret request"))
		return
	}

	// Create the secret metadata
	meta := &SecretMetadata{
		Password: req.Password, // TODO: argon2 hash the password since there is no reason for us to store raw passwords
		Filename: req.Filename,
		IsBase64: req.IsBase64,
		Created:  time.Now(),
	}

	// Compute the number of accesses for the secret
	if req.Accesses == 0 {
		meta.Accesses = DefaultSecretAccesses
	} else {
		meta.Accesses = req.Accesses
	}

	// Compute the expiration time from the request
	if req.Lifetime == v1.Duration(0) {
		meta.Expires = meta.Created.Add(DefaultSecretLifetime)
	} else {
		meta.Expires = meta.Created.Add(time.Duration(req.Lifetime))
	}

	// Create the reply back to the user
	var err error
	rep := &v1.CreateSecretReply{
		Expires: meta.Expires,
	}

	// Make a random URL to store the secret in
	if rep.Token, err = GenerateUniqueURL(); err != nil {
		log.Error().Err(err).Msg("could not generate unique token for secret")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// Store the password in the database
	tmpSecretsStore[rep.Token] = req.Secret
	tmpSecretsMeta[rep.Token] = meta

	// Return success
	c.JSON(http.StatusCreated, rep)
}

// FetchSecret handles an incoming fetch secret request and attempts to retrieve the
// secret from the database and return it to the user. This function also handles the
// password and ensures that a 404 is returned to obfuscate the existence of the secret
// on bad requests.
// TODO: handle passwords with argon2
func (s *Server) FetchSecret(c *gin.Context) {
	// Fetch the meta with the token
	token := c.Param("token")
	meta, ok := tmpSecretsMeta[token]
	if !ok {
		c.JSON(http.StatusNotFound, ErrorResponse("secret not found"))
		return
	}

	// Check the secret is valid prior to returning a response (in case a sidechannel
	// retrieval or race condition failed to destroy the password).
	if !meta.Valid() {
		log.Warn().Msg("race condition or invalid secret metadata fetched, destroying")
		delete(tmpSecretsMeta, token)
		delete(tmpSecretsStore, token)
		c.JSON(http.StatusNotFound, ErrorResponse("secret not found"))
		return
	}

	// Check the password if it has been posted
	// TODO: perform argon2 password checking
	if meta.Password != "" {
		// A password is required as an Authorization: Bearer <token> header where the
		// token is the base64 encoded password. Basic auth does not apply here since
		// there is no username associated with the secret.
		password := ParseBearerToken(c.GetHeader("Authorization"))
		if password == "" || password != meta.Password {
			c.JSON(http.StatusUnauthorized, ErrorResponse("password required for secret"))
			return
		}
	}

	// Update metadata with the access info
	meta.Access()

	// Create the secret reply
	rep := v1.FetchSecretReply{
		Filename: meta.Filename,
		IsBase64: meta.IsBase64,
		Created:  meta.Created,
		Accesses: meta.Retrievals,
	}

	// Fetch the secret with the token
	if rep.Secret, ok = tmpSecretsStore[token]; !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse("unhandled secret store error"))
		return
	}

	// Cleanup if necessary
	if !meta.Valid() {
		delete(tmpSecretsMeta, token)
		delete(tmpSecretsStore, token)
	}

	// Return the successful reply
	c.JSON(http.StatusOK, rep)
}

// DestroySecret handles an incoming destroy secret request and attempts to delete the
// secret from the database. This RPC is password protected in the same way fetch is.
// TODO: handle passwords with argon2
func (s *Server) DestroySecret(c *gin.Context) {
	// Fetch the meta with the token
	token := c.Param("token")
	meta, ok := tmpSecretsMeta[token]
	if !ok {
		c.JSON(http.StatusNotFound, ErrorResponse("secret not found"))
		return
	}

	// Check the secret is valid prior to returning a response (in case a sidechannel
	// retrieval or race condition failed to destroy the password).
	if !meta.Valid() {
		log.Warn().Msg("race condition or invalid secret metadata fetched, destroying")
		delete(tmpSecretsMeta, token)
		delete(tmpSecretsStore, token)
		c.JSON(http.StatusNotFound, ErrorResponse("secret not found"))
		return
	}

	// Check the password if it has been posted
	// TODO: perform argon2 password checking
	if meta.Password != "" {
		// A password is required as an Authorization: Bearer <token> header where the
		// token is the base64 encoded password. Basic auth does not apply here since
		// there is no username associated with the secret.
		password := ParseBearerToken(c.GetHeader("Authorization"))
		if password == "" || password != meta.Password {
			c.JSON(http.StatusUnauthorized, ErrorResponse("password required for secret"))
			return
		}
	}

	// Delete the secret from the database
	delete(tmpSecretsMeta, token)
	delete(tmpSecretsStore, token)

	// Return the successful reply
	c.JSON(http.StatusOK, &v1.DestroySecretReply{Destroyed: true})
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

var bearerRegex = regexp.MustCompile(`^(?i)Bearer\s+([A-Za-z0-9=+/_-]+)$`)

// ParseBearerToken parses an Authorization: Bearer <token> header such that the token
// is the base64 encoded password. Basic auth not used here since there is no user.
func ParseBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header != "" && bearerRegex.MatchString(header) {
		groups := bearerRegex.FindStringSubmatch(header)
		token, err := base64.URLEncoding.DecodeString(groups[1])
		if err != nil {
			log.Warn().Err(err).Msg("could not base64 decode Authorization header")
			return ""
		}
		return string(token)
	}
	return ""
}
