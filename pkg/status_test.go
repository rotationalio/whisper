package whisper_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	whisper "github.com/rotationalio/whisper/pkg"
)

func (s *WhisperTestSuite) TestStatus() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/status", nil)
	s.router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	s.Equal(http.StatusOK, result.StatusCode)
	s.Equal("application/json; charset=utf-8", result.Header.Get("Content-Type"))

	var data map[string]interface{}
	err := json.NewDecoder(result.Body).Decode(&data)
	s.NoError(err)

	s.Contains(data, "status")
	s.Equal("ok", data["status"])
	s.Contains(data, "version")
	s.Equal(whisper.Version(), data["version"])

}
