package vault

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	smpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go"
	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewMock creates and returns a client to access a mock Secret Manager for testing.
// Note that the SecretManager is identical and all external functionality is unchanged,
// however instead of making requests to Google Secret Manager, the mock object is
// simply storing things in memory.
func NewMock(conf config.GoogleConfig) (*SecretManager, error) {
	return &SecretManager{
		parent: fmt.Sprintf("projects/%s", conf.Project),
		client: &mockSecretManagerClient{
			secrets: make(map[string]*mockSecret),
		},
	}, nil
}

type mockSecretManagerClient struct {
	secrets map[string]*mockSecret
}

type mockSecret struct {
	Name     string    // looks like [parent]/secrets/[token]-[suffix]
	Created  time.Time // time the secreted was created
	Expires  time.Time // when the secret expires
	Versions [][]byte  // secret versions contain data
}

func (c *mockSecretManagerClient) GetSecret(ctx context.Context, req *smpb.GetSecretRequest, opts ...gax.CallOption) (*smpb.Secret, error) {
	log.Warn().Str("method", "GetSecret").Msg("mock secret manager called")
	// Check if secret is in the mock database
	if secret, ok := c.secrets[req.Name]; ok && secret.Expires.After(time.Now()) {
		return &smpb.Secret{
			Name:       secret.Name,
			CreateTime: timestamppb.New(secret.Created),
		}, nil
	}

	// Otherwise return not found
	return nil, status.Error(codes.NotFound, "secret not found")
}

func (c *mockSecretManagerClient) CreateSecret(ctx context.Context, req *smpb.CreateSecretRequest, opts ...gax.CallOption) (*smpb.Secret, error) {
	log.Warn().Str("method", "CreateSecret").Msg("mock secret manager called")
	if req.Parent == "" || req.SecretId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing parent or secret id")
	}

	secret := &mockSecret{
		Name:     fmt.Sprintf("%s/secrets/%s", req.Parent, req.SecretId),
		Created:  time.Now(),
		Versions: make([][]byte, 0),
	}

	// Handle the expiration
	switch expires := req.Secret.Expiration.(type) {
	case *smpb.Secret_ExpireTime:
		secret.Expires = expires.ExpireTime.AsTime()
		if secret.Expires.IsZero() {
			return nil, status.Error(codes.InvalidArgument, "invalid expiration time")
		}
	case *smpb.Secret_Ttl:
		ttl := expires.Ttl.AsDuration()
		if ttl < 1 {
			return nil, status.Error(codes.InvalidArgument, "invalid time to live")
		}
		secret.Expires = secret.Created.Add(ttl)
	default:
		return nil, status.Error(codes.InvalidArgument, "unknown expiration type")
	}

	// Check if secret already exists
	if _, ok := c.secrets[secret.Name]; ok {
		return nil, status.Error(codes.AlreadyExists, "secret already exists")
	}

	// Add secret to the "database"
	c.secrets[secret.Name] = secret

	return &smpb.Secret{
		Name:       secret.Name,
		CreateTime: timestamppb.New(secret.Created),
	}, nil
}

func (c *mockSecretManagerClient) AddSecretVersion(ctx context.Context, req *smpb.AddSecretVersionRequest, opts ...gax.CallOption) (*smpb.SecretVersion, error) {
	log.Warn().Str("method", "AddSecretVersion").Msg("mock secret manager called")
	if req.Parent == "" {
		return nil, status.Error(codes.InvalidArgument, "missing parent")
	}

	if len(req.Payload.Data) > 66560 {
		return nil, status.Error(codes.InvalidArgument, "payload too large")
	}

	secret, ok := c.secrets[req.Parent]
	if !ok {
		return nil, status.Error(codes.NotFound, "secret not found")
	}

	if secret.Expires.Before(time.Now()) {
		delete(c.secrets, req.Parent)
		return nil, status.Error(codes.NotFound, "secret expired")
	}

	// Add the version to the database and return the version
	// TODO: do we need to populate any of the other secret version fields?
	secret.Versions = append(secret.Versions, req.Payload.Data)
	return &smpb.SecretVersion{
		Name: fmt.Sprintf("%s/versions/%d", secret.Name, len(secret.Versions)),
	}, nil
}

func (c *mockSecretManagerClient) AccessSecretVersion(ctx context.Context, req *smpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*smpb.AccessSecretVersionResponse, error) {
	log.Warn().Str("method", "AccessSecretVersion").Msg("mock secret manager called")
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "missing secret version name")
	}

	parts := strings.Split(req.Name, "/")
	if len(parts) < 3 || parts[len(parts)-2] != "versions" {
		return nil, status.Error(codes.InvalidArgument, "could not parse secret version name")
	}
	parent := strings.Join(parts[:len(parts)-2], "/")

	secret, ok := c.secrets[parent]
	if !ok {
		return nil, status.Error(codes.NotFound, "secret not found")
	}

	if secret.Expires.Before(time.Now()) {
		delete(c.secrets, parent)
		return nil, status.Error(codes.NotFound, "secret expired")
	}

	var idx int64
	if strings.ToLower(parts[len(parts)-1]) == "latest" {
		idx = int64(len(secret.Versions) - 1)
	} else {
		var err error
		if idx, err = strconv.ParseInt(parts[len(parts)-1], 10, 0); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	if idx >= int64(len(secret.Versions)) {
		return nil, status.Error(codes.NotFound, "version not found")
	}

	return &smpb.AccessSecretVersionResponse{
		Name: fmt.Sprintf("%s/versions/%d", secret.Name, idx+1),
		Payload: &smpb.SecretPayload{
			Data: secret.Versions[idx],
		},
	}, nil
}

func (c *mockSecretManagerClient) DeleteSecret(ctx context.Context, req *smpb.DeleteSecretRequest, opts ...gax.CallOption) error {
	log.Warn().Str("method", "DeleteSecret").Msg("mock secret manager called")
	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "missing secret name")
	}

	secret, ok := c.secrets[req.Name]
	if !ok {
		return status.Error(codes.NotFound, "secret not found")
	}

	delete(c.secrets, req.Name)
	if secret.Expires.Before(time.Now()) {
		return status.Error(codes.NotFound, "secret expired")
	}
	return nil
}
