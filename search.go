package gondole

import (
	"encoding/json"
	"fmt"

	"github.com/sendgrid/rest"
)

// Search search for contents (accounts or statuses) and returns a Results
func (g *Client) Search(query string, resolve bool) (*Results, error) {
	if query == "" {
		return nil, ErrInvalidParameter
	}
	req := g.prepareRequest("search")
	req.QueryParams["q"] = query
	if resolve {
		req.QueryParams["resolve"] = "true"
	}
	r, err := rest.API(req)
	if err != nil {
		return nil, fmt.Errorf("search: %s", err.Error())
	}

	// Check for error reply
	var errorResult Error
	if err := json.Unmarshal([]byte(r.Body), &errorResult); err == nil {
		// The empty object is not an error
		if errorResult.Text != "" {
			return nil, fmt.Errorf("%s", errorResult.Text)
		}
	}

	// Not an error reply; let's unmarshal the data
	var results Results
	err = json.Unmarshal([]byte(r.Body), &results)
	if err != nil {
		return nil, fmt.Errorf("search API: %s", err.Error())
	}
	return &results, nil
}
