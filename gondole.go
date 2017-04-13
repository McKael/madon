package gondole

import (
	"errors"
	"fmt"

	"github.com/sendgrid/rest"
)

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

// prepareRequest insert all pre-defined stuff
func (g *Client) prepareRequest(what string) (req rest.Request) {
	var endPoint string

	endPoint = g.APIBase + "/" + what
	// Add at least one option, the APIkey if present
	hdrs := make(map[string]string)
	opts := make(map[string]string)

	// Insert our sig
	hdrs["User-Agent"] = fmt.Sprintf("Gondole/%s", GondoleVersion)
	if g.userToken != nil {
		hdrs["Authorization"] = fmt.Sprintf("Bearer %s", g.userToken.AccessToken)
	}

	req = rest.Request{
		BaseURL:     endPoint,
		Headers:     hdrs,
		QueryParams: opts,
	}
	return
}
