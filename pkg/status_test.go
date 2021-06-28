package whisper_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	whisper "github.com/rotationalio/whisper/pkg"
	"github.com/stretchr/testify/require"
)

func (s *WhisperTestSuite) TestStatus() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/status", nil)
	s.router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	require.Equal(s.T(), http.StatusOK, w.Code)
	require.Equal(s.T(), "application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	require.NoError(s.T(), err)

	require.Contains(s.T(), data, "status")
	require.Equal(s.T(), "ok", data["status"])
	require.Contains(s.T(), data, "version")
	require.Equal(s.T(), whisper.Version(), data["version"])

}
