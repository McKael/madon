package gondole

import ()
import (
	"github.com/sendgrid/rest"
	"fmt"
)

const (
	APIVersion = "0.0"

	APIEndpoint = "/api/v1"

	NoRedirect = "urn:ietf:wg:oauth:2.0:oob"
)

// prepareRequest insert all pre-defined stuff
func (g *Gondole) prepareRequest(what string) (req rest.Request) {
	endPoint := APIEndpoint + fmt.Sprintf("/%s/", what)
	key, ok := HasAPIKey()

	// Add at least one option, the APIkey if present
	hdrs := make(map[string]string)
	opts := make(map[string]string)

	// Insert our sig
	hdrs["User-Agent"] = fmt.Sprintf("Gondole/%s", APIVersion)

	// Insert key
	if ok {
		opts["key"] = key
	}

	req = rest.Request{
		BaseURL:     endPoint,
		Headers:     hdrs,
		QueryParams: opts,
	}
	return
}
