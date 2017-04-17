/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package gondole

import (
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// getAccountsOptions contains option fields for POST and DELETE API calls
type getAccountsOptions struct {
	// The ID is used for most commands
	ID int

	// The following fields are used when searching for accounts
	Q     string
	Limit int
}

// getSingleAccount returns an account entity
// The operation 'op' can be "account", "verify_credentials", "follow",
// "unfollow", "block", "unblock", "mute", "unmute",
// "follow_requests/authorize" or // "follow_requests/reject".
// The id is optional and depends on the operation.
func (g *Client) getSingleAccount(op string, id int) (*Account, error) {
	var endPoint string
	method := rest.Get
	strID := strconv.Itoa(id)

	switch op {
	case "account":
		endPoint = "accounts/" + strID
	case "verify_credentials":
		endPoint = "accounts/verify_credentials"
	case "follow", "unfollow", "block", "unblock", "mute", "unmute":
		endPoint = "accounts/" + strID + "/" + op
		method = rest.Post
	case "follow_requests/authorize", "follow_requests/reject":
		// The documentation is incorrect, the endpoint actually
		// is "follow_requests/:id/{authorize|reject}"
		endPoint = op[:16] + strID + "/" + op[16:]
		method = rest.Post
	default:
		return nil, ErrInvalidParameter
	}

	var account Account
	if err := g.apiCall(endPoint, method, nil, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// getMultipleAccounts returns a list of account entities
// The operation 'op' can be "followers", "following", "search", "blocks",
// "mutes", "follow_requests".
// The id is optional and depends on the operation.
func (g *Client) getMultipleAccounts(op string, opts *getAccountsOptions) ([]Account, error) {
	var endPoint string

	switch op {
	case "followers", "following":
		if opts == nil || opts.ID < 1 {
			return []Account{}, ErrInvalidID
		}
		endPoint = "accounts/" + strconv.Itoa(opts.ID) + "/" + op
	case "follow_requests", "blocks", "mutes":
		endPoint = op
	case "search":
		if opts == nil || opts.Q == "" {
			return []Account{}, ErrInvalidParameter
		}
		endPoint = "accounts/" + op
	default:
		return nil, ErrInvalidParameter
	}

	// Handle target-specific query parameters
	params := make(apiCallParams)
	if op == "search" {
		params["q"] = opts.Q
		if opts.Limit > 0 {
			params["limit"] = strconv.Itoa(opts.Limit)
		}
	}

	var accounts []Account
	if err := g.apiCall(endPoint, rest.Get, params, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// GetAccount returns an account entity
// The returned value can be nil if there is an error or if the
// requested ID does not exist.
func (g *Client) GetAccount(accountID int) (*Account, error) {
	account, err := g.getSingleAccount("account", accountID)
	if err != nil {
		return nil, err
	}
	if account != nil && account.ID == 0 {
		return nil, ErrEntityNotFound
	}
	return account, nil
}

// GetCurrentAccount returns the current user account
func (g *Client) GetCurrentAccount() (*Account, error) {
	account, err := g.getSingleAccount("verify_credentials", 0)
	if err != nil {
		return nil, err
	}
	if account != nil && account.ID == 0 {
		return nil, ErrEntityNotFound
	}
	return account, nil
}

// GetAccountFollowers returns the list of accounts following a given account
func (g *Client) GetAccountFollowers(accountID int) ([]Account, error) {
	o := &getAccountsOptions{ID: accountID}
	return g.getMultipleAccounts("followers", o)
}

// GetAccountFollowing returns the list of accounts a given account is following
func (g *Client) GetAccountFollowing(accountID int) ([]Account, error) {
	o := &getAccountsOptions{ID: accountID}
	return g.getMultipleAccounts("following", o)
}

// FollowAccount follows an account
func (g *Client) FollowAccount(accountID int) error {
	account, err := g.getSingleAccount("follow", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// UnfollowAccount unfollows an account
func (g *Client) UnfollowAccount(accountID int) error {
	account, err := g.getSingleAccount("unfollow", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// FollowRemoteAccount follows a remote account
// The parameter 'uri' is a URI (e.g. "username@domain").
func (g *Client) FollowRemoteAccount(uri string) (*Account, error) {
	if uri == "" {
		return nil, ErrInvalidID
	}

	params := make(apiCallParams)
	params["uri"] = uri

	var account Account
	if err := g.apiCall("follows", rest.Post, params, &account); err != nil {
		return nil, err
	}
	if account.ID == 0 {
		return nil, ErrEntityNotFound
	}
	return &account, nil
}

// BlockAccount blocks an account
func (g *Client) BlockAccount(accountID int) error {
	account, err := g.getSingleAccount("block", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// UnblockAccount unblocks an account
func (g *Client) UnblockAccount(accountID int) error {
	account, err := g.getSingleAccount("unblock", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// MuteAccount mutes an account
func (g *Client) MuteAccount(accountID int) error {
	account, err := g.getSingleAccount("mute", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// UnmuteAccount unmutes an account
func (g *Client) UnmuteAccount(accountID int) error {
	account, err := g.getSingleAccount("unmute", accountID)
	if err != nil {
		return err
	}
	if account != nil && account.ID != accountID {
		return ErrEntityNotFound
	}
	return nil
}

// SearchAccounts returns a list of accounts matching the query string
// The limit parameter is optional (can be 0).
func (g *Client) SearchAccounts(query string, limit int) ([]Account, error) {
	o := &getAccountsOptions{Q: query, Limit: limit}
	return g.getMultipleAccounts("search", o)
}

// GetBlockedAccounts returns the list of blocked accounts
func (g *Client) GetBlockedAccounts() ([]Account, error) {
	return g.getMultipleAccounts("blocks", nil)
}

// GetMutedAccounts returns the list of muted accounts
func (g *Client) GetMutedAccounts() ([]Account, error) {
	return g.getMultipleAccounts("mutes", nil)
}

// GetAccountFollowRequests returns the list of follow requests accounts
func (g *Client) GetAccountFollowRequests() ([]Account, error) {
	return g.getMultipleAccounts("follow_requests", nil)
}

// GetAccountRelationships returns a list of relationship entities for the given accounts
func (g *Client) GetAccountRelationships(accountIDs []int) ([]Relationship, error) {
	if len(accountIDs) < 1 {
		return nil, ErrInvalidID
	}

	params := make(apiCallParams)
	for i, id := range accountIDs {
		if id < 1 {
			return nil, ErrInvalidID
		}
		qID := fmt.Sprintf("id[%d]", i+1)
		params[qID] = strconv.Itoa(id)
	}

	var rl []Relationship
	if err := g.apiCall("accounts/relationships", rest.Get, params, &rl); err != nil {
		return nil, err
	}
	return rl, nil
}

// GetAccountStatuses returns a list of status entities for the given account
// If onlyMedia is true, returns only statuses that have media attachments.
// If excludeReplies is true, skip statuses that reply to other statuses.
func (g *Client) GetAccountStatuses(accountID int, onlyMedia, excludeReplies bool) ([]Status, error) {
	if accountID < 1 {
		return nil, ErrInvalidID
	}

	endPoint := "accounts/" + strconv.Itoa(accountID) + "/" + "statuses"
	params := make(apiCallParams)
	if onlyMedia {
		params["only_media"] = "true"
	}
	if excludeReplies {
		params["exclude_replies"] = "true"
	}

	var sl []Status
	if err := g.apiCall(endPoint, rest.Get, params, &sl); err != nil {
		return nil, err
	}
	return sl, nil
}

// FollowRequestAuthorize authorizes or rejects an account follow-request
func (g *Client) FollowRequestAuthorize(accountID int, authorize bool) error {
	endPoint := "follow_requests/reject"
	if authorize {
		endPoint = "follow_requests/authorize"
	}
	_, err := g.getSingleAccount(endPoint, accountID)
	return err
}
