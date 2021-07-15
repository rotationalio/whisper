package whisper

import (
	"context"
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

// CreateSecret handles an incoming CreateSecretRequest and attempts to create a new
// secret that will only be displayed when the correct link is retrieved.
func (s *Server) CreateSecret(c *gin.Context) {
	// Parse incoming JSON data from the client request
	var req v1.CreateSecretRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Warn().Err(err).Msg("could not bind request")
		c.JSON(http.StatusBadRequest, ErrorResponse("invalid create secret request"))
		return
	}

	// Make a random URL to store the secret in
	var (
		err   error
		token string
	)
	if token, err = s.GenerateUniqueURL(context.TODO()); err != nil {
		log.Error().Err(err).Msg("could not generate unique token for secret")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// Create the secret context
	meta := s.vault.With(token)
	meta.Filename = req.Filename
	meta.IsBase64 = req.IsBase64
	meta.Created = time.Now()

	// Store the password as a derived key
	if err = meta.SetPassword(req.Password); err != nil {
		log.Error().Err(err).Msg("could not create derived key")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	// Compute the number of accesses for the secret
	if req.Accesses == 0 {
		meta.Accesses = DefaultSecretAccesses
		log.Debug().Int("accesses", meta.Accesses).Msg("using default number of accesses")
	} else {
		meta.Accesses = req.Accesses
		log.Debug().Int("accesses", meta.Accesses).Msg("using user supplied number of accesses")
	}

	// Compute the expiration time from the request
	if req.Lifetime == v1.Duration(0) {
		meta.Expires = meta.Created.Add(DefaultSecretLifetime)
		log.Debug().Dur("ttl", DefaultSecretLifetime).Msg("using default secret lifetime")
	} else {
		meta.Expires = meta.Created.Add(time.Duration(req.Lifetime))
		log.Debug().Dur("ttl", time.Duration(req.Lifetime)).Msg("using user supplied secret lifetime")
	}

	// Create the secret in the vault.
	if err = meta.New(context.TODO(), req.Secret); err != nil {
		log.Error().Err(err).Msg("could not create new secret in vault")
		c.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}

	// Return successful reply back to the user
	c.JSON(http.StatusCreated, &v1.CreateSecretReply{
		Token:   token,
		Expires: meta.Expires,
	})
}

// FetchSecret handles an incoming fetch secret request and attempts to retrieve the
// secret from the database and return it to the user. This function also handles the
// password and ensures that a 404 is returned to obfuscate the existence of the secret
// on bad requests.
func (s *Server) FetchSecret(c *gin.Context) {
	// Prepare to fetch the meta with the token and password from the request
	token := c.Param("token")
	meta := s.vault.With(token)
	password := ParseBearerToken(c.GetHeader("Authorization"))
	log.Debug().Bool("authorization", password != "").Msg("beginning fetch")

	// Attempt to retrieve the secret from the database
	secret, err := meta.Fetch(context.TODO(), password)
	if err != nil {
		switch err {
		case ErrSecretNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse(err))
		case ErrNotAuthorized:
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
		default:
			log.Error().Err(err).Msg("could not fetch secret")
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		}
		return
	}

	// Create the secret reply
	rep := v1.FetchSecretReply{
		Secret:   secret,
		Filename: meta.Filename,
		IsBase64: meta.IsBase64,
		Created:  meta.Created,
		Accesses: meta.Retrievals,
	}

	// Return the successful reply
	c.JSON(http.StatusOK, rep)
}

// DestroySecret handles an incoming destroy secret request and attempts to delete the
// secret from the database. This RPC is password protected in the same way fetch is.
func (s *Server) DestroySecret(c *gin.Context) {
	// Prepare to fetch the meta with the token and password from the request
	token := c.Param("token")
	meta := s.vault.With(token)
	password := ParseBearerToken(c.GetHeader("Authorization"))
	log.Debug().Bool("authorization", password != "").Msg("beginning destroy")

	// Delete the secret from the database
	// Attempt to retrieve the secret from the database
	err := meta.Destroy(context.TODO(), password)
	if err != nil {
		switch err {
		case ErrSecretNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse(err))
		case ErrNotAuthorized:
			c.JSON(http.StatusUnauthorized, ErrorResponse(err))
		default:
			log.Error().Err(err).Msg("could not destroy secret")
			c.JSON(http.StatusInternalServerError, ErrorResponse(err))
		}
		return
	}

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
func (s *Server) GenerateUniqueURL(ctx context.Context) (token string, err error) {
	for i := 0; i < generateUniqueAttempts; i++ {
		// Create a random array of bytes
		buf := make([]byte, generateUniqueLength)
		if _, err = rand.Read(buf); err != nil {
			return "", err
		}

		// Convert random array into URL safe base64 encoded string
		token := base64.RawURLEncoding.EncodeToString(buf)

		// Check if the token exists already in the database
		var exists bool
		if exists, err = s.vault.Check(ctx, token); err != nil {
			return "", err
		}

		if !exists {
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
