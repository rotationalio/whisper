package whisper_test

import (
	"net/http"
	"net/http/httptest"
)

func (s *WhisperTestSuite) TestVersionRedirect() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.router.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	s.Equal(http.StatusPermanentRedirect, result.StatusCode)
	s.Equal("application/json; charset=utf-8", result.Header.Get("Content-Type"))
	s.Equal("/v1", result.Header.Get("Location"))
}
