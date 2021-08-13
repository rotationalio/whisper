package api_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rotationalio/whisper/pkg/api/v1"
	"github.com/stretchr/testify/require"
)

func TestDuration(t *testing.T) {
	// Test marshal duration
	duration := api.Duration(time.Millisecond * 1586)
	data, err := json.Marshal(duration)
	require.NoError(t, err)
	require.Equal(t, "\"1.586s\"", string(data))

	// Test unmarshal from string
	var dur api.Duration
	require.NoError(t, json.Unmarshal([]byte("\"1m32s\""), &dur))
	require.Equal(t, api.Duration(time.Second*92), dur)

	// Test unmarshal from float
	var dur2 api.Duration
	require.NoError(t, json.Unmarshal([]byte("5255000000"), &dur2))
	require.Equal(t, api.Duration(time.Millisecond*5255), dur2)
}
