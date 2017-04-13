package gondole

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// updateStatusOptions contains option field for POST and DELETE API calls
type updateStatusOptions struct {
	// The ID is used for most commands
	ID int

	// The following fields are used for posting a new status
	Status      string
	InReplyToID int
	MediaIDs    []int
	Sensitive   bool
	SpoilerText string
	Visibility  string // "direct", "private", "unlisted" or "public"
}

// queryStatusData queries the statuses API
// The subquery can be empty or "status" (the status itself), "context",
// "card", "reblogged_by", "favourited_by".
// The data argument will receive the object(s) returned by the API server.
func (g *Client) queryStatusData(statusID int, subquery string, data interface{}) error {
	endPoint := "statuses/" + strconv.Itoa(statusID)

	if statusID < 1 {
		return ErrInvalidID
	}

	if subquery != "" && subquery != "status" {
		switch subquery {
		case "context", "card", "reblogged_by", "favourited_by":
		default:
			return ErrInvalidParameter
		}

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

// updateStatusData updates the statuses
// The subquery can be empty or "status" (to post a status), "delete" (for
// deleting a status), "reblog", "unreblog", "favourite", "unfavourite".
// The data argument will receive the object(s) returned by the API server.
func (g *Client) updateStatusData(subquery string, opts updateStatusOptions, data interface{}) error {
	method := rest.Post
	endPoint := "statuses"

	switch subquery {
	case "", "status":
		subquery = "status"
		if opts.Status == "" {
			return ErrInvalidParameter
		}
		switch opts.Visibility {
		case "", "direct", "private", "unlisted", "public":
			// Okay
		default:
			return ErrInvalidParameter
		}
		if len(opts.MediaIDs) > 4 {
			return fmt.Errorf("too many (>4) media IDs")
		}
	case "delete":
		method = rest.Delete
		if opts.ID < 1 {
			return ErrInvalidID
		}
		endPoint += "/" + strconv.Itoa(opts.ID)
	case "reblog", "unreblog", "favourite", "unfavourite":
		if opts.ID < 1 {
			return ErrInvalidID
		}
		endPoint += "/" + strconv.Itoa(opts.ID) + "/" + subquery
	default:
		return ErrInvalidParameter
	}

	req := g.prepareRequest(endPoint)
	req.Method = method

	// Form items for a new toot
	if subquery == "status" {
		req.QueryParams["status"] = opts.Status
		if opts.InReplyToID > 0 {
			req.QueryParams["in_reply_to_id"] = strconv.Itoa(opts.InReplyToID)
		}
		for i, id := range opts.MediaIDs {
			qID := fmt.Sprintf("media_ids[%d]", i+1)
			req.QueryParams[qID] = strconv.Itoa(id)
		}
		if opts.Sensitive {
			req.QueryParams["sensitive"] = "true"
		}
		if opts.SpoilerText != "" {
			req.QueryParams["spoiler_text"] = opts.SpoilerText
		}
		if opts.Visibility != "" {
			req.QueryParams["visibility"] = opts.Visibility
		}
	}

	r, err := rest.API(req)
	if err != nil {
		return fmt.Errorf("status/%s API query: %s", subquery, err.Error())
	}

	// Check for error reply
	var errorResult Error
	if err := json.Unmarshal([]byte(r.Body), &errorResult); err == nil {
		return fmt.Errorf("%s", errorResult.Text)
	}

	// Not an error reply; let's unmarshall the data
	err = json.Unmarshal([]byte(r.Body), &data)
	if err != nil {
		return fmt.Errorf("status/%s API: %s", subquery, err.Error())
	}
	return nil
}

// GetStatus returns a status
// The returned status can be nil if there is an error or if the
// requested ID does not exist.
func (g *Client) GetStatus(id int) (*Status, error) {
	var status Status

	if err := g.queryStatusData(id, "status", &status); err != nil {
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

// PostStatus posts a new "toot"
// All parameters but "text" can be empty.
// Visibility must be empty, or one of "direct", "private", "unlisted" and "public".
func (g *Client) PostStatus(text string, inReplyTo int, mediaIDs []int, sensitive bool, spoilerText string, visibility string) (*Status, error) {
	var status Status
	o := updateStatusOptions{
		Status:      text,
		InReplyToID: inReplyTo,
		MediaIDs:    mediaIDs,
		Sensitive:   sensitive,
		SpoilerText: spoilerText,
		Visibility:  visibility,
	}

	err := g.updateStatusData("status", o, &status)
	if err != nil {
		return nil, err
	}
	if status.ID == 0 {
		return nil, ErrEntityNotFound // TODO Change error message
	}
	return &status, err
}

// DeleteStatus deletes a status
func (g *Client) DeleteStatus(id int) error {
	var status Status
	o := updateStatusOptions{ID: id}
	err := g.updateStatusData("delete", o, &status)
	return err
}

// ReblogStatus reblogs a status
func (g *Client) ReblogStatus(id int) error {
	var status Status
	o := updateStatusOptions{ID: id}
	err := g.updateStatusData("reblog", o, &status)
	return err
}

// UnreblogStatus unreblogs a status
func (g *Client) UnreblogStatus(id int) error {
	var status Status
	o := updateStatusOptions{ID: id}
	err := g.updateStatusData("unreblog", o, &status)
	return err
}

// FavouriteStatus favourites a status
func (g *Client) FavouriteStatus(id int) error {
	var status Status
	o := updateStatusOptions{ID: id}
	err := g.updateStatusData("favourite", o, &status)
	return err
}

// UnfavouriteStatus unfavourites a status
func (g *Client) UnfavouriteStatus(id int) error {
	var status Status
	o := updateStatusOptions{ID: id}
	err := g.updateStatusData("unfavourite", o, &status)
	return err
}
