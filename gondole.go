package gondole

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sendgrid/rest"
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

// prepareRequest inserts all pre-defined stuff
func (g *Client) prepareRequest(target string, method rest.Method, params apiCallParams) (req rest.Request) {
	endPoint := g.APIBase + "/" + target

	// Request headers
	hdrs := make(map[string]string)
	hdrs["User-Agent"] = fmt.Sprintf("Gondole/%s", GondoleVersion)
	if g.userToken != nil {
		hdrs["Authorization"] = fmt.Sprintf("Bearer %s", g.userToken.AccessToken)
	}

	req = rest.Request{
		BaseURL:     endPoint,
		Headers:     hdrs,
		Method:      method,
		QueryParams: params,
	}
	return
}

// apiCall makes a call to the Mastodon API server
func (g *Client) apiCall(endPoint string, method rest.Method, params apiCallParams, data interface{}) error {
	// Prepare query
	req := g.prepareRequest(endPoint, method, params)

	// Make API call
	r, err := rest.API(req)
	if err != nil {
		return fmt.Errorf("API query (%s) failed: %s", endPoint, err.Error())
	}

	// Check for error reply
	var errorResult Error
	if err := json.Unmarshal([]byte(r.Body), &errorResult); err == nil {
		// The empty object is not an error
		if errorResult.Text != "" {
			return fmt.Errorf("%s", errorResult.Text)
		}
	}

	// Not an error reply; let's unmarshal the data
	err = json.Unmarshal([]byte(r.Body), &data)
	if err != nil {
		return fmt.Errorf("cannot decode API response (%s): %s", method, err.Error())
	}
	return nil
}
