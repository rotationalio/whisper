package whisper

import "time"

// SecretMetadata stores sidechannel information related to the secret but not the
// secret itself. This data allows the whipser service to manage passwords, the number
// of accesses, and the expiration of the secret without having to retrieve the secret
// directly, creating a possible vulnerability.
type SecretMetadata struct {
	Password     string    `json:"password,omitempty"` // the argon2 hashed password for comparision
	Filename     string    `json:"filename,omitempty"` // if the secret is a file, the name of the file for download
	IsBase64     bool      `json:"is_base64"`          // if the secret is base64 encoded or not
	Accesses     int       `json:"accesses"`           // the number of allowed accesses for the secret
	Retrievals   int       `json:"retrievals"`         // counts the number of times the secret has been accessed
	Created      time.Time `json:"created"`            // the timestamp the secret was created
	LastAccessed time.Time `json:"last_accessed"`      // the timestamp that the secret was last accessed
	Expires      time.Time `json:"expires"`            // the timestamp when the secret will have expired
}

// Valid returns true if the retrievals is less than the number of allowed accesses and
// the current time is before the expiration time.
func (s *SecretMetadata) Valid() bool {
	if time.Now().After(s.Expires) {
		return false
	}

	if s.Accesses > 0 && s.Retrievals >= s.Accesses {
		return false
	}

	return true
}

// Access updates the secret metadata on a fetch or other access to the secret.
func (s *SecretMetadata) Access() {
	s.Retrievals++
	s.LastAccessed = time.Now()
}
