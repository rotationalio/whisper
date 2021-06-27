package whisper_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/rotationalio/whisper/pkg"
	"github.com/stretchr/testify/require"
)

func TestNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	NotFound(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusNotFound)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)

}

func TestNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	NotAllowed(ctx)

	result := r.Result()
	defer result.Body.Close()
	require.Equal(t, result.StatusCode, http.StatusMethodNotAllowed)
	require.Equal(t, "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(t, err)
}
