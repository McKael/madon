/*
Copyright 2017-2018 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"github.com/sendgrid/rest"
)

// Search search for contents (accounts or statuses) and returns a Results
func (mc *Client) Search(query string, resolve bool) (*Results, error) {
	if query == "" {
		return nil, ErrInvalidParameter
	}

	params := make(apiCallParams)
	params["q"] = query
	if resolve {
		params["resolve"] = "true"
	}

	var resultsV1 struct {
		Results
		Hashtags []string `json:"hashtags"`
	}
	if err := mc.apiCall("search", rest.Get, params, nil, nil, &resultsV1); err != nil {
		return nil, err
	}

	var results Results
	results.Accounts = resultsV1.Accounts
	results.Statuses = resultsV1.Statuses
	for _, t := range resultsV1.Hashtags {
		results.Hashtags = append(results.Hashtags, Tag{Name: t})
	}

	return &results, nil
}
