package vault_test

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/rotationalio/whisper/pkg/config"
	"github.com/rotationalio/whisper/pkg/vault"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VaultTestSuite struct {
	suite.Suite
	conf  config.GoogleConfig
	vault *vault.SecretManager
}

func (s *VaultTestSuite) SetupSuite() {
	zerolog.SetGlobalLevel(zerolog.PanicLevel)

	s.conf = config.GoogleConfig{
		Credentials: "path/to/nowhere.json",
		Project:     "vault-test-project",
	}

	var err error
	s.vault, err = vault.NewMock(s.conf)
	s.NoError(err)
}

func TestVault(t *testing.T) {
	suite.Run(t, new(VaultTestSuite))
}

func (s *VaultTestSuite) TestSecretContextFlow() {
	// Create the secret manager context
	token := createToken()
	secret := s.vault.With(token)

	// Update the metadata
	secret.Filename = ""
	secret.IsBase64 = false
	secret.Accesses = 1
	secret.Created = time.Now()
	secret.Expires = time.Now().Add(24 * time.Hour)

	// Create the secret and metadata
	s.NoError(secret.New(context.TODO(), "the eagle flies at midnight"))
	s.Empty(secret.Password)

	// Attempt to fetch the secret
	secret2 := s.vault.With(token)
	whisper, destroyed, err := secret2.Fetch(context.TODO(), "")
	s.NoError(err)
	s.True(destroyed)
	s.Equal("the eagle flies at midnight", whisper)

	// Now that the secret has been destroyed fetch should return not found
	secret3 := s.vault.With(token)
	_, _, err = secret3.Fetch(context.TODO(), "")
	s.ErrorIs(err, vault.ErrSecretNotFound)
}

func (s *VaultTestSuite) TestSecretContextPasswordFlow() {
	// Create the secret manager context
	token := createToken()
	secret := s.vault.With(token)

	// Update the metadata
	secret.Filename = ""
	secret.IsBase64 = false
	secret.Accesses = 1
	secret.Created = time.Now()
	secret.Expires = time.Now().Add(24 * time.Hour)

	secret.SetPassword("theunlock")
	s.NotEmpty(secret.Password)

	// Create the secret and metadata
	s.NoError(secret.New(context.TODO(), "the eagle flies at midnight"))

	// Attempt to fetch the secret without a password
	secret2 := s.vault.With(token)
	_, _, err := secret2.Fetch(context.TODO(), "")
	s.ErrorIs(err, vault.ErrNotAuthorized)

	// Attempt with the wrong password
	secret3 := s.vault.With(token)
	_, _, err = secret3.Fetch(context.TODO(), "opensaysme")
	s.ErrorIs(err, vault.ErrNotAuthorized)

	// Finally, make sure the password can be retrieved
	secret4 := s.vault.With(token)
	whisper, destroyed, err := secret4.Fetch(context.TODO(), "theunlock")
	s.NoError(err)
	s.True(destroyed)
	s.Equal("the eagle flies at midnight", whisper)

	// Now that the secret has been destroyed fetch should return not found
	secret5 := s.vault.With(token)
	_, _, err = secret5.Fetch(context.TODO(), "theunlock")
	s.ErrorIs(err, vault.ErrSecretNotFound)
}

func (s *VaultTestSuite) TestDestroy() {
	// Create the secret manager context
	token := createToken()
	secret := s.vault.With(token)

	// Update the metadata
	secret.Filename = ""
	secret.IsBase64 = false
	secret.Accesses = 1
	secret.Created = time.Now()
	secret.Expires = time.Now().Add(24 * time.Hour)

	// Create the secret and metadata
	s.NoError(secret.New(context.TODO(), "the eagle flies at midnight"))
	s.Empty(secret.Password)

	// Destroy the secret
	secret2 := s.vault.With(token)
	s.NoError(secret2.Destroy(context.TODO(), ""))

	// Now that the secret has been destroyed fetch should return not found
	secret3 := s.vault.With(token)
	_, _, err := secret3.Fetch(context.TODO(), "theunlock")
	s.ErrorIs(err, vault.ErrSecretNotFound)
}

func (s *VaultTestSuite) TestDestroyPassword() {
	// Create the secret manager context
	token := createToken()
	secret := s.vault.With(token)

	// Update the metadata
	secret.Filename = ""
	secret.IsBase64 = false
	secret.Accesses = 1
	secret.Created = time.Now()
	secret.Expires = time.Now().Add(24 * time.Hour)

	secret.SetPassword("theunlock")
	s.NotEmpty(secret.Password)

	// Create the secret and metadata
	s.NoError(secret.New(context.TODO(), "the eagle flies at midnight"))

	// Attempt to fetch the secret without a password
	secret2 := s.vault.With(token)
	err := secret2.Destroy(context.TODO(), "")
	s.ErrorIs(err, vault.ErrNotAuthorized)

	// Attempt with the wrong password
	secret3 := s.vault.With(token)
	err = secret3.Destroy(context.TODO(), "opensaysme")
	s.ErrorIs(err, vault.ErrNotAuthorized)

	// Finally, make sure the password can be deleted
	secret4 := s.vault.With(token)
	s.NoError(secret4.Destroy(context.TODO(), "theunlock"))

	// Now that the secret has been destroyed fetch should return not found
	secret5 := s.vault.With(token)
	_, _, err = secret5.Fetch(context.TODO(), "theunlock")
	s.ErrorIs(err, vault.ErrSecretNotFound)
}

func (s *VaultTestSuite) TestCheckEmpty() {
	token := createToken()
	found, err := s.vault.Check(context.TODO(), token)
	s.NoError(err)
	s.False(found)
}

func (s *VaultTestSuite) TestSecretContextMethods() {
	// No manager or token, secret is invalid
	secret := &vault.SecretContext{}
	s.False(secret.Valid())

	// Has manager and token, but still invalid without timestamps
	secret = s.vault.With(createToken())
	s.False(secret.Valid())
	secret.Created = time.Now()
	s.False(secret.Valid())

	// Although it now has expires, it's after now so invalid
	secret.Expires = time.Now().Add(-24 * time.Hour)
	s.False(secret.Valid())

	// Even without accesses or retrievals, should now be valid
	secret.Expires = time.Now().Add(24 * time.Hour)
	s.True(secret.Valid())

	// It is valid to have a secret context without a last accessed timestamp
	s.Zero(secret.LastAccessed)

	// Accesses is now greter than 0, but retrievals is still 0
	secret.Accesses = 2
	s.True(secret.Valid())

	// One retrieval keeps the context valid
	secret.Access()
	s.True(secret.Valid())

	// After too many accesses, the secret becomes invalid
	secret.Access()
	s.False(secret.Valid())

	// Access should have added a timestamp
	s.NotZero(secret.LastAccessed)
}

func (s *VaultTestSuite) TestSecretContextPassword() {
	secret := s.vault.With(createToken())
	s.Empty(secret.Password)

	// Set the password
	s.NoError(secret.SetPassword("zoolander"))
	s.NotEmpty(secret.Password)

	// Cannot verify the password unless it has been loaded from the vault
	s.ErrorIs(secret.VerifyPassword("zoolander"), vault.ErrNotLoaded)

	// Should be able to remove the password
	s.NoError(secret.SetPassword(""))
	s.Empty(secret.Password)

	// Cannot verify the password unless it has been loaded from the vault, even if empty
	s.ErrorIs(secret.VerifyPassword(""), vault.ErrNotLoaded)
}

func (s *VaultTestSuite) TestExpiresRequired() {
	secret := s.vault.With(createToken())
	err := secret.New(context.TODO(), "the first secret should be completed easily")
	s.ErrorIs(err, vault.ErrTimeToLive)
}

func (s *VaultTestSuite) TestDuplicateNew() {
	token := createToken()
	secret := s.vault.With(token)
	secret.Expires = time.Now().Add(24 * time.Hour)
	s.NoError(secret.New(context.TODO(), "the first secret should be completed easily"))

	secret2 := s.vault.With(token)
	secret2.Expires = time.Now().Add(32 * time.Hour)
	err := secret2.New(context.TODO(), "a different secret with the same token")
	s.ErrorIs(err, vault.ErrAlreadyExists)
}

func TestCreateToken(t *testing.T) {
	tokens := make([]string, 0, 10)
	for i := 0; i < 10; i++ {
		tokens = append(tokens, createToken())
	}

	for i, token := range tokens {
		require.Len(t, token, 43)
		for j, other := range tokens {
			if i == j {
				continue
			}
			require.NotEqual(t, token, other)
		}
	}
}

func createToken() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}
