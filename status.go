package gondole

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// queryStatusData queries the statuses API
// The subquery can be empty (the status itself), "context", "card",
// "reblogged_by", "favourited_by".
func (g *Client) queryStatusData(statusID int, subquery string, data interface{}) error {
	endPoint := "statuses/" + strconv.Itoa(statusID)

	if statusID < 1 {
		return ErrInvalidID
	}

	if subquery != "" {
		// TODO: check subquery values?
		endPoint += "/" + subquery
	}
	req := g.prepareRequest(endPoint)
	r, err := rest.API(req)
	if err != nil {
		return fmt.Errorf("status/%s API query: %s", subquery, err.Error())
	}

	err = json.Unmarshal([]byte(r.Body), &data)
	if err != nil {
		var errorRes Error
		err2 := json.Unmarshal([]byte(r.Body), &errorRes)
		if err2 == nil {
			return fmt.Errorf("%s", errorRes.Text)
		}
		return fmt.Errorf("status/%s API: %s", subquery, err.Error())
	}

	return nil
}

// GetStatus returns a status
// The returned status can be nil if there is an error or if the
// requested ID does not exist.
func (g *Client) GetStatus(id int) (*Status, error) {
	var status Status

	if err := g.queryStatusData(id, "", &status); err != nil {
		return nil, err
	}

	if status.ID == 0 {
		return nil, ErrEntityNotFound
	}

	return &status, nil
}

// GetStatusContext returns a status context
func (g *Client) GetStatusContext(id int) (*Context, error) {
	var context Context

	if err := g.queryStatusData(id, "context", &context); err != nil {
		return nil, err
	}

	return &context, nil
}

// GetStatusCard returns a status card
func (g *Client) GetStatusCard(id int) (*Card, error) {
	var card Card

	if err := g.queryStatusData(id, "card", &card); err != nil {
		return nil, err
	}

	return &card, nil
}

// GetStatusRebloggedBy returns a list of the accounts who reblogged a status
func (g *Client) GetStatusRebloggedBy(id int) ([]Account, error) {
	var accounts []Account
	err := g.queryStatusData(id, "reblogged_by", &accounts)
	return accounts, err
}

// GetStatusFavouritedBy returns a list of the accounts who favourited a status
func (g *Client) GetStatusFavouritedBy(id int) ([]Account, error) {
	var accounts []Account
	err := g.queryStatusData(id, "favourited_by", &accounts)
	return accounts, err
}
