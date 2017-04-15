package gondole

import (
	"errors"
)

// apiCallParams is a map with the parameters for an API call
type apiCallParams map[string]string

const (
	// GondoleVersion contains the version of the Gondole implementation
	GondoleVersion = "0.0"

	defaultInstanceURL = "https://mastodon.social"
	apiVersion         = "v1" // That is not overridable
	defaultAPIPath     = "/api/" + apiVersion

	// NoRedirect is the URI for no redirection in the App registration
	NoRedirect = "urn:ietf:wg:oauth:2.0:oob"
)

// Error codes
var (
	ErrAlreadyRegistered = errors.New("app already registered")
	ErrEntityNotFound    = errors.New("entity not found")
	ErrInvalidParameter  = errors.New("incorrect parameter")
	ErrInvalidID         = errors.New("incorrect entity ID")
)
