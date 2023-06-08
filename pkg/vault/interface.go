package vault

import (
	"context"

	smpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/googleapis/gax-go"
)

// secretManagerClient describes the methods used to interact with the Google Secret
// Manager, primarily to allow mocking this interface for testing purposes. It is also
// conceivable that this interface could be used to define other vault storage.
type secretManagerClient interface {
	GetSecret(ctx context.Context, req *smpb.GetSecretRequest, opts ...gax.CallOption) (*smpb.Secret, error)
	CreateSecret(ctx context.Context, req *smpb.CreateSecretRequest, opts ...gax.CallOption) (*smpb.Secret, error)
	AddSecretVersion(ctx context.Context, req *smpb.AddSecretVersionRequest, opts ...gax.CallOption) (*smpb.SecretVersion, error)
	AccessSecretVersion(ctx context.Context, req *smpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*smpb.AccessSecretVersionResponse, error)
	DeleteSecret(ctx context.Context, req *smpb.DeleteSecretRequest, opts ...gax.CallOption) error
}
