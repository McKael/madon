/*
Copyright 2017 Ollivier Robert
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package gondole

import (
	"net/url"
	"strings"

	"github.com/sendgrid/rest"
)

type registerApp struct {
	ID           int    `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NewApp registers a new instance
func NewApp(name string, scopes []string, redirectURI, instanceURL string) (g *Client, err error) {
	if instanceURL == "" {
		instanceURL = defaultInstanceURL
	}

	if !strings.Contains(instanceURL, "://") {
		instanceURL = "https://" + instanceURL
	}

	apiPath := instanceURL + defaultAPIPath

	if _, err := url.ParseRequestURI(apiPath); err != nil {
		return nil, err
	}

	g = &Client{
		Name:        name,
		APIBase:     apiPath,
		InstanceURL: instanceURL,
	}

	params := make(apiCallParams)
	params["client_name"] = name
	params["scopes"] = strings.Join(scopes, " ")
	if redirectURI != "" {
		params["redirect_uris"] = redirectURI
	} else {
		params["redirect_uris"] = NoRedirect
	}

	var app registerApp
	if err := g.apiCall("apps", rest.Post, params, &app); err != nil {
		return nil, err
	}

	g.ID = app.ClientID
	g.Secret = app.ClientSecret

	return
}
