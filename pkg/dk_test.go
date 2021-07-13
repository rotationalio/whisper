package whisper_test

import (
	"testing"

	. "github.com/rotationalio/whisper/pkg"
	"github.com/stretchr/testify/require"
)

func TestDerivedKey(t *testing.T) {
	// Create a derived key from a password
	passwd, err := CreateDerivedKey("theeaglefliesatmidnight")
	require.NoError(t, err)

	verified, err := VerifyDerivedKey(passwd, "theeaglefliesatmidnight")
	require.NoError(t, err)
	require.True(t, verified)

	verified, err = VerifyDerivedKey(passwd, "thesearentthedroidsyourelookingfor")
	require.NoError(t, err)
	require.False(t, verified)

	// Create a derived key from a password
	passwd2, err := CreateDerivedKey("lightning")
	require.NoError(t, err)
	require.NotEqual(t, passwd, passwd2)
}
