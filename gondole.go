package gondole

import (
	"github.com/sendgrid/rest"
	"fmt"
	"errors"
)

const (
	APIVersion = "0.0"

	// That is not overridable
	apiURL = "/api/v1"

	// Fallback instance
	FallBackURL = "https://mastodon.social"

	NoRedirect = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	APIEndpoint string
	ErrAlreadyRegistered = errors.New("App already registered")
)

// prepareRequest insert all pre-defined stuff
func (g *Client) prepareRequest(what string) (req rest.Request) {
	var APIBase string

	// Allow for overriding for registration
	if g.APIBase == "" {
		APIBase = FallBackURL
	} else {
		APIBase = g.APIBase
	}

	APIEndpoint = fmt.Sprintf("%s%s/%", APIBase, apiURL, what)

	// Add at least one option, the APIkey if present
	hdrs := make(map[string]string)
	opts := make(map[string]string)

	// Insert our sig
	hdrs["User-Agent"] = fmt.Sprintf("Client/%s", APIVersion)
	hdrs["Authorization"] = fmt.Sprintf("Bearer %s", g.Secret)

	req = rest.Request{
		BaseURL:     APIEndpoint,
		Headers:     hdrs,
		QueryParams: opts,
	}
	return
}

func (g *Client) Login() (err error) {
	return
}
