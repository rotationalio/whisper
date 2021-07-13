package whisper_test

import (
	"encoding/base64"
	"testing"

	. "github.com/rotationalio/whisper/pkg"
	"github.com/stretchr/testify/require"
)

func TestGenerateUniqueURL(t *testing.T) {
	tokens := make(map[string]struct{})
	for i := 0; i < 48; i++ {
		// Generate token
		token, err := GenerateUniqueURL()
		require.NoError(t, err)
		require.Len(t, token, 43)

		// Make sure token is unique
		_, ok := tokens[token]
		require.False(t, ok)

		// Add token to unique set
		tokens[token] = struct{}{}
	}
}

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
