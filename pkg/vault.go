package whisper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rs/zerolog/log"
	smpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Suffixes describe the name that the secret or metadata is stored with
const (
	SuffixSecret   = "secret"
	SuffixMetadata = "metadata"
)

// Standard errors for error type checking
var (
	ErrAlreadyExists    = errors.New("secret already exists")
	ErrSecretNotFound   = errors.New("secret does not exist in secret manager")
	ErrFileSizeLimit    = errors.New("secret payload exceeds size limit")
	ErrPermissionDenied = errors.New("secret manager permission denied")
	ErrNotAuthorized    = errors.New("correct password required")
	ErrNotLoaded        = errors.New("secret context needs to be loaded")
)

// NewSecretManager creates and returns a client to access the Google Secret Manager.
// This function requires the $GOOGLE_APPLICATION_CREDENTIALS environment variable to
// be set, which specifies the JSON path to the service account credentials.
func NewSecretManager(conf config.GoogleConfig) (sm *SecretManager, err error) {
	sm = &SecretManager{parent: fmt.Sprintf("projects/%s", conf.Project)}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if sm.client, err = secretmanager.NewClient(ctx); err != nil {
		return nil, fmt.Errorf("could not connect to secret manager: %s", err)
	}

	return sm, nil
}

// SecretManager provides access to the Google Secret Manager and is the primary "vault"
// (secret storage) currently used by Whisper. The manager maintains the secret parent
// path composed by the project name as well as the RPC client.
type SecretManager struct {
	parent string
	client *secretmanager.Client
}

// With extracts a secret context with the information required to fetch a secret from
// Google Secret Manager. This is used to create a new context and to retrieve one.
func (sm *SecretManager) With(token string) *SecretContext {
	return &SecretContext{
		manager: sm,
		token:   token,
	}
}

// Check returns true if the secret exists, false if it does not. Used to determine if
// the secret exists as quickly as possible (e.g. to ensure no duplicates).
func (sm *SecretManager) Check(ctx context.Context, token string) (_ bool, err error) {
	// Build the request to add the version based on the standardized path, using metadata suffix.
	path := fmt.Sprintf("%s/secrets/%s-%s", sm.parent, token, SuffixMetadata)
	req := &smpb.GetSecretRequest{
		Name: path,
	}

	// Create an internal context to avoid an infinite hang by a failed API call
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the request
	if _, err = sm.client.GetSecret(sctx, req); err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return false, err
		}

		// If this is not a context error, attempt to parse the gRPC status error
		serr, ok := status.FromError(err)
		if ok {
			// Log the original message since it will be subsumed by the error check
			log.Debug().Err(err).Msg("get secret rpc error")

			switch serr.Code() {
			case codes.NotFound:
				// If the secret doesn't exist (e.g. not created yet or deleted)
				// This is the condition we're looking for, so no error.
				return false, nil
			case codes.PermissionDenied:
				// If we've given a wrong path, wrong project, or wrong service account
				return false, ErrPermissionDenied
			}
		}

		// If the error is something else, something went wrong.
		return false, fmt.Errorf("could not get secret: %s", err)
	}

	// If there was no error, we assume that we retrieved the secret.
	return true, nil
}

// SecretContext stores sidechannel information related to the secret but not the
// secret itself. This data allows the whipser service to manage passwords, the number
// of accesses, and the expiration of the secret without having to retrieve the secret
// directly, creating a possible vulnerability. The context is also responsible for
// managing interactions with the Google Secret Manager service for a specific secret,
// including using the derived key algorithm for password verification and checking.
type SecretContext struct {
	// External information that is serialized and stored in the secret manager.
	Password     string    `json:"password,omitempty"` // the argon2 hashed password for comparision
	Filename     string    `json:"filename,omitempty"` // if the secret is a file, the name of the file for download
	IsBase64     bool      `json:"is_base64"`          // if the secret is base64 encoded or not
	Accesses     int       `json:"accesses"`           // the number of allowed accesses for the secret
	Retrievals   int       `json:"retrievals"`         // counts the number of times the secret has been accessed
	Created      time.Time `json:"created"`            // the timestamp the secret was created
	LastAccessed time.Time `json:"last_accessed"`      // the timestamp that the secret was last accessed
	Expires      time.Time `json:"expires"`            // the timestamp when the secret will have expired

	// Internal information required to access secret manager api.
	manager *SecretManager // client to make calls to the service
	token   string         // the token that the context is stored with
	loaded  bool           // if the context has been loaded from the database or not
}

// SetPassword is the preferred way for setting a password on a secret that is about to
// be created since it guarantees that the derived key methodology is correct.
func (s *SecretContext) SetPassword(password string) (err error) {
	// If no password is supplied then this secret does not require a password.
	// NOTE: we could create a dervied key from an empty password, but this would
	// increase the time it would take to retrieve a password, and it's not clear if
	// that is valuable in the case where there is no password.
	if password == "" {
		s.Password = ""
		return nil
	}

	// Otherwise create the derived key for the password.
	if s.Password, err = CreateDerivedKey(password); err != nil {
		return err
	}
	return nil
}

// Valid returns true if the retrievals is less than the number of allowed accesses and
// the current time is before the expiration time. If the Expires or Created timestamp
// is zero, the context is assumed to not have been initialized. Valid is used both to
// check if the secret context can be created/updated and to determine if it should be
// destroyed.
func (s *SecretContext) Valid() bool {
	if s.manager == nil || s.token == "" {
		return false
	}

	if s.Expires.IsZero() || s.Created.IsZero() {
		return false
	}

	if time.Now().After(s.Expires) {
		return false
	}

	if s.Accesses > 0 && s.Retrievals >= s.Accesses {
		return false
	}

	return true
}

// Access updates the secret metadata on a fetch or other access to the secret.
func (s *SecretContext) Access() {
	s.Retrievals++
	s.LastAccessed = time.Now()
}

// New creates a new secret and metadata in Google Secret Manager adding the first
// version to actually store the data. Returns an error if the secret already exists.
func (s *SecretContext) New(ctx context.Context, secret string) (err error) {
	// Marshal the context first so that if anything goes wrong we don't strand data in
	// Google Secret Manager.
	var data []byte
	if data, err = json.Marshal(s); err != nil {
		return fmt.Errorf("could not marshal secret metadata: %s", err)
	}

	// Create the metadata secret
	if err = s.Create(ctx, SuffixMetadata); err != nil {
		return fmt.Errorf("could not create metadata: %s", err)
	}

	// Add a version for the metadata
	if err = s.AddVersion(ctx, SuffixMetadata, data); err != nil {
		log.Warn().Bool("metadata", true).Bool("secret", false).Msg("incomplete secret creation")
		return fmt.Errorf("could not add metadata version: %s", err)
	}

	// Create the secret next
	if err = s.Create(ctx, SuffixSecret); err != nil {
		log.Warn().Bool("metadata version", true).Bool("secret", false).Msg("incomplete secret creation")
		return fmt.Errorf("could not create secret actual: %s", err)
	}

	// Add a version for the secret
	if err = s.AddVersion(ctx, SuffixSecret, []byte(secret)); err != nil {
		log.Warn().Bool("metadata version", true).Bool("secret", true).Msg("incomplete secret creation")
		return fmt.Errorf("could not add secret actual version: %s", err)
	}

	return nil
}

// Fetch loads the metadata into the context, then determines if a password is required
// and validates the password using the derived key algorithm. If the secret metadata is
// still valid then it returns the secret, otherwise it returns not found. The metadata
// can still be accessed if it was loaded from the Secret Manager.
func (s *SecretContext) Fetch(ctx context.Context, password string) (_ string, err error) {
	// First fetch the secret metadata
	if err = s.Load(ctx, false); err != nil {
		return "", err
	}

	// Check the secret is valid prior to returning a response (in case a sidechannel
	// retrieval or race condition failed to destroy the password).
	if !s.Valid() {
		log.Warn().Msg("race condition or invalid secret metadata fetched, destroying")
		if err = s.Destroy(ctx, password); err != nil {
			log.Error().Err(err).Msg("could not destroy invalid secret")
		}
		return "", ErrSecretNotFound
	}

	// Check if the password is required and if so, if it matches the derived key.
	if err = s.VerifyPassword(password); err != nil {
		return "", err
	}

	// Fetch the latest version of the secret
	var secret []byte
	if secret, err = s.LatestVersion(ctx, SuffixSecret); err != nil {
		return "", err
	}

	// Update the metadata with the access information
	s.Access()

	// Store the updated metadata or cleanup if necessary
	if s.Valid() {
		var payload []byte
		if payload, err = json.Marshal(s); err != nil {
			return "", fmt.Errorf("could not marshal secret context: %s", err)
		}

		// Update the metadata with the new version
		if err = s.AddVersion(ctx, SuffixMetadata, payload); err != nil {
			return "", fmt.Errorf("could not update metadata: %s", err)
		}
	} else {
		// Don't return the error in this case because the secret will eventually expire
		log.Debug().Msg("destroying now invalid secret after access")
		if err = s.Destroy(ctx, password); err != nil {
			log.Error().Err(err).Msg("could not destroy invalid secret after access")
		}
	}

	return string(secret), nil
}

// Destroy both the secret metadata and the secret unless the password is incorrect
// (returns not authorized) or the secret does not exist (returns not found).
func (s *SecretContext) Destroy(ctx context.Context, password string) (err error) {
	// First load the secret metadata - won't load if already loaded.
	if err = s.Load(ctx, false); err != nil {
		return err
	}

	// Check if the password is required and if so, if it matches the derived key.
	// Otherwise anyone could destroy a secret. This only matters if the secret is still
	// valid, if it's not valid; destroy no matter what the password is.
	if s.Valid() {
		if err = s.VerifyPassword(password); err != nil {
			return err
		}
	}

	// Delete the secret first
	if err = s.Delete(ctx, SuffixSecret); err != nil {
		return fmt.Errorf("could not delete secret actual: %s", err)
	}

	// Delete the metadata last
	if err = s.Delete(ctx, SuffixMetadata); err != nil {
		return fmt.Errorf("could not delete secret metadata: %s", err)
	}

	return nil
}

// Load is a helper function that retrieves the secret metadata from the Secret Manager.
// It is safe to call load multiple times because it will only load once unless reload
func (s *SecretContext) Load(ctx context.Context, reload bool) (err error) {
	// Check if the Secret has been loaded already (and we're not reloading)
	if s.loaded && !reload {
		return nil
	}

	// Fetch the secret metadata from Secret Manager
	var payload []byte
	if payload, err = s.LatestVersion(ctx, SuffixMetadata); err != nil {
		// LatestVersion will return the error not found if necessary
		return err
	}

	// Parse the payload into the context
	if err = json.Unmarshal(payload, s); err != nil {
		return fmt.Errorf("could not unmarshal secret metadata: %s", err)
	}

	s.loaded = true
	return nil
}

// Create is an helper function that is called twice from New: once to create the secret
// metadata and once to create the secret itself. The only external information required
// is the token which is stored on the context.
func (s *SecretContext) Create(ctx context.Context, suffix string) (err error) {
	// Build the request to create the secret in the specified parent where the ID is
	// the token + suffix (e.g. token-secret or token-metadata).
	req := &smpb.CreateSecretRequest{
		Parent:   s.manager.parent,
		SecretId: fmt.Sprintf("%s-%s", s.token, suffix),
		Secret: &smpb.Secret{
			Expiration: &smpb.Secret_ExpireTime{
				ExpireTime: timestamppb.New(s.Expires),
			},
			Replication: &smpb.Replication{
				Replication: &smpb.Replication_Automatic_{
					Automatic: &smpb.Replication_Automatic{},
				},
			},
		},
	}

	// Create an internal context, since a failed API call will result in infinite hang
	// Note that the outer context is the parent of the subcontext.
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Call the API. Note: We don't actually need the result that comes back from the API call
	// and not accessing it directly (e.g. logging plaintext, etc) provides added security
	if _, err = s.manager.client.CreateSecret(sctx, req); err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// If the secret already exists, return an error that can be checked
		serr, ok := status.FromError(err)
		if ok && serr.Code() == codes.AlreadyExists {
			log.Debug().Err(err).Msg("create secret rpc error")
			return ErrAlreadyExists
		}

		// If the error is something else, something went wrong.
		return fmt.Errorf("could not create secret: [%d] %s", serr.Code(), serr.Message())
	}
	return nil
}

// AddVersion updates the Secret with the new payload and is a helper function that is
// used both in New to create the first version and in Fetch to track accesses and
// updates in the secret metadata.
func (s *SecretContext) AddVersion(ctx context.Context, suffix string, payload []byte) (err error) {
	// Build the request to add the version based on the standardized path and suffix.
	path := fmt.Sprintf("%s/secrets/%s-%s", s.manager.parent, s.token, suffix)
	req := &smpb.AddSecretVersionRequest{
		Parent: path,
		Payload: &smpb.SecretPayload{
			Data: payload,
		},
	}

	// Create an internal context to avoid an infinite hang by a failed API call
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the request. Note: we don't actually need the result from the API call
	// and we're not accessing it directly to ensure we don't leak sensitive info.
	if _, err = s.manager.client.AddSecretVersion(sctx, req); err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// If this is not a context error, attempt to parse the gRPC status error
		serr, ok := status.FromError(err)
		if ok {
			log.Debug().Err(err).Msg("add secret version rpc error")
			switch serr.Code() {
			case codes.NotFound:
				// If the secret doesn't exist (e.g. not created yet or deleted)
				return ErrSecretNotFound
			case codes.InvalidArgument:
				// Maximum size limit of 65KiB for the payload
				return ErrFileSizeLimit
			case codes.PermissionDenied:
				// If we've given a wrong path, wrong project, or wrong service account
				return ErrPermissionDenied
			}
		}

		// If the error is something else, something went wrong.
		return fmt.Errorf("could not add %q version: %s", suffix, err)
	}

	return nil
}

// LatestVersion returns the payload for the latest version of the secret if it exists.
// This is a helper function that performs no validation or password verification.
func (s *SecretContext) LatestVersion(ctx context.Context, suffix string) (_ []byte, err error) {
	// Build the request to add the version based on the standardized path and suffix.
	path := fmt.Sprintf("%s/secrets/%s-%s/versions/latest", s.manager.parent, s.token, suffix)
	req := &smpb.AccessSecretVersionRequest{
		Name: path,
	}

	// Create an internal context to avoid an infinite hang by a failed API call
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the request
	var result *smpb.AccessSecretVersionResponse
	if result, err = s.manager.client.AccessSecretVersion(sctx, req); err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}

		// If this is not a context error, attempt to parse the gRPC status error
		serr, ok := status.FromError(err)
		if ok {
			log.Debug().Err(err).Msg("access secret version rpc error")
			switch serr.Code() {
			case codes.NotFound:
				// If the secret doesn't exist (e.g. not created yet or deleted)
				return nil, ErrSecretNotFound
			case codes.PermissionDenied:
				// If we've given a wrong path, wrong project, or wrong service account
				return nil, ErrPermissionDenied
			}
		}

		// If the error is something else, something went wrong.
		return nil, fmt.Errorf("could not fetch %q latest version: %s", suffix, err)
	}

	return result.Payload.Data, nil
}

func (s *SecretContext) Delete(ctx context.Context, suffix string) (err error) {
	// Build the request to add the version based on the standardized path and suffix.
	path := fmt.Sprintf("%s/secrets/%s-%s", s.manager.parent, s.token, suffix)
	req := &smpb.DeleteSecretRequest{
		Name: path,
	}

	// Create an internal context to avoid an infinite hang by a failed API call
	sctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Execute the request
	if err = s.manager.client.DeleteSecret(sctx, req); err != nil {
		// If the API call is malformed, it will hang until the internal context times out
		if errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		// If this is not a context error, attempt to parse the gRPC status error
		serr, ok := status.FromError(err)
		if ok {
			log.Debug().Err(err).Msg("delete secret rpc error")
			switch serr.Code() {
			case codes.NotFound:
				// If the secret doesn't exist (e.g. not created yet or deleted)
				return ErrSecretNotFound
			case codes.PermissionDenied:
				// If we've given a wrong path, wrong project, or wrong service account
				return ErrPermissionDenied
			}
		}

		// If the error is something else, something went wrong.
		return fmt.Errorf("could not delete %q secret: %s", suffix, err)
	}

	return nil
}

// VerifyPassword checks that the password matches the dervied password otherwise errors.
func (s *SecretContext) VerifyPassword(password string) (err error) {
	if !s.loaded {
		return ErrNotLoaded
	}

	if s.Password != "" {
		if password == "" {
			log.Debug().Msg("password required but no password supplied")
			return ErrNotAuthorized
		}

		var verified bool
		if verified, err = VerifyDerivedKey(s.Password, password); err != nil {
			return err
		}
		if !verified {
			log.Debug().Msg("incorrect password supplied")
			return ErrNotAuthorized
		}
	}
	return nil
}
