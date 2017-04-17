/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package gondole

import (
	"github.com/sendgrid/rest"
)

// GetCurrentInstance returns current instance information
func (g *Client) GetCurrentInstance() (*Instance, error) {
	var i Instance
	if err := g.apiCall("instance", rest.Get, nil, &i); err != nil {
		return nil, err
	}
	return &i, nil
}
