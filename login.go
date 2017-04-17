/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package gondole

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sendgrid/rest"
)

// UserToken represents a user token as returned by the Mastodon API
type UserToken struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

// LoginBasic does basic user authentication
func (g *Client) LoginBasic(username, password string, scopes []string) error {
	if g == nil {
		return fmt.Errorf("use of uninitialized gondole client")
	}

	if username == "" {
		return fmt.Errorf("missing username")
	}
	if password == "" {
		return fmt.Errorf("missing password")
	}

	hdrs := make(map[string]string)
	opts := make(map[string]string)

	hdrs["User-Agent"] = "Gondole/" + GondoleVersion

	opts["grant_type"] = "password"
	opts["client_id"] = g.ID
	opts["client_secret"] = g.Secret
	opts["username"] = username
	opts["password"] = password
	if len(scopes) > 0 {
		opts["scope"] = strings.Join(scopes, " ")
	}

	req := rest.Request{
		BaseURL:     g.InstanceURL + "/oauth/token",
		Headers:     hdrs,
		QueryParams: opts,
		Method:      rest.Post,
	}

	r, err := restAPI(req)
	if err != nil {
		return err
	}

	var resp UserToken

	err = json.Unmarshal([]byte(r.Body), &resp)
	if err != nil {
		return fmt.Errorf("cannot unmarshal server response: %s", err.Error())
	}

	g.UserToken = &resp
	return nil
}

// SetUserToken sets an existing user credentials
// No verification of the arguments is made.
func (g *Client) SetUserToken(token, username, password string, scopes []string) error {
	if g == nil {
		return fmt.Errorf("use of uninitialized gondole client")
	}

	g.UserToken = &UserToken{
		AccessToken: token,
		Scope:       strings.Join(scopes, " "),
		TokenType:   "bearer",
	}
	return nil
}
