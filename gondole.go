package gondole

import (
	"github.com/sendgrid/rest"
	"fmt"
	"errors"
)

const (
	APIVersion = "0.0"

	defAPIEndpoint = "https://mastodon.social/api/v1"

	NoRedirect = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	APIEndpoint string
	ErrAlreadyRegistered = errors.New("App already registered")
)

// prepareRequest insert all pre-defined stuff
func (g *Gondole) prepareRequest(what string) (req rest.Request) {
	var endPoint string

	// Allow for overriding for registration
	if APIEndpoint == "" {
		endPoint = defAPIEndpoint
	} else {
		endPoint = APIEndpoint
	}

	endPoint = endPoint + fmt.Sprintf("/%s", what)
	// Add at least one option, the APIkey if present
	hdrs := make(map[string]string)
	opts := make(map[string]string)

	// Insert our sig
	hdrs["User-Agent"] = fmt.Sprintf("Gondole/%s", APIVersion)
	hdrs["Authorization"] = fmt.Sprintf("Bearer %s", g.Secret)

	req = rest.Request{
		BaseURL:     endPoint,
		Headers:     hdrs,
		QueryParams: opts,
	}
	return
}
