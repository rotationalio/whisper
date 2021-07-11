package whisper_test

import (
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
		require.Len(t, token, 44)

		// Make sure token is unique
		_, ok := tokens[token]
		require.False(t, ok)

		// Add token to unique set
		tokens[token] = struct{}{}
	}
}
