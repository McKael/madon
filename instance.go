package gondole

import (
	"encoding/json"
	"fmt"

	"github.com/sendgrid/rest"
)

// GetCurrentInstance returns current instance information
func (g *Client) GetCurrentInstance() (*Instance, error) {
	req := g.prepareRequest("instance")
	r, err := rest.API(req)
	if err != nil {
		return nil, fmt.Errorf("instance: %s", err.Error())
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
	var i Instance
	err = json.Unmarshal([]byte(r.Body), &i)
	if err != nil {
		return nil, fmt.Errorf("instance API: %s", err.Error())
	}
	return &i, nil
}
