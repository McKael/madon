package gondole

import (
	"github.com/sendgrid/rest"
)

// GetFavourites returns the list of the user's favourites
func (g *Client) GetFavourites() ([]Status, error) {
	var faves []Status
	err := g.apiCall("favourites", rest.Get, nil, &faves)
	if err != nil {
		return nil, err
	}
	return faves, nil
}
