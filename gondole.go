package gondole

import (
	"github.com/sendgrid/rest"
	"fmt"
	"errors"
)

const (
	APIVersion = "0.0"

	APIEndpoint = "/api/v1"

	NoRedirect = "urn:ietf:wg:oauth:2.0:oob"
)

var (
	ErrAlreadyRegistered = errors.New("App already registered")
)

// prepareRequest insert all pre-defined stuff
func (g *Gondole) prepareRequest(what string) (req rest.Request) {
	endPoint := APIEndpoint + fmt.Sprintf("/%s/", what)

	// Add at least one option, the APIkey if present
	hdrs := make(map[string]string)
	opts := make(map[string]string)

	// Insert our sig
	hdrs["User-Agent"] = fmt.Sprintf("Gondole/%s", APIVersion)

	req = rest.Request{
		BaseURL:     endPoint,
		Headers:     hdrs,
		QueryParams: opts,
	}
	return
}
