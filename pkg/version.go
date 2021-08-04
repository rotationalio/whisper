package whisper

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Version component constants for the current build.
const (
	VersionMajor         = 1
	VersionMinor         = 0
	VersionPatch         = 1
	VersionReleaseLevel  = ""
	VersionReleaseNumber = 0
)

// Version returns the semantic version for the current build.
func Version() string {
	var versionCore string
	if VersionPatch > 0 {
		versionCore = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	} else {
		versionCore = fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
	}

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			return fmt.Sprintf("%s-%s%d", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		}
		return fmt.Sprintf("%s-%s", versionCore, VersionReleaseLevel)
	}
	return versionCore
}

// VersionURL returns the URL prefix for the API at the current version
func VersionURL() string {
	return fmt.Sprintf("/v%d", VersionMajor)
}

// RedirectVersion sends the caller to the root of the current version
func (s *Server) RedirectVersion(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Redirect(http.StatusPermanentRedirect, VersionURL())
}
