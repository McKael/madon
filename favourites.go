package gondole

import (
	"encoding/json"
	"fmt"

	"github.com/sendgrid/rest"
)

// GetFavourites returns the list of the user's favourites
func (g *Client) GetFavourites() ([]Status, error) {
	var faves []Status

	req := g.prepareRequest("favourites")
	r, err := rest.API(req)
	if err != nil {
		return faves, fmt.Errorf("favourites API query: %s", err.Error())
	}

	println(r.Body)
	err = json.Unmarshal([]byte(r.Body), &faves)
	if err != nil {
		var errorRes Error
		err2 := json.Unmarshal([]byte(r.Body), &errorRes)
		if err2 == nil {
			return faves, fmt.Errorf("%s", errorRes.Text)
		}
		return faves, fmt.Errorf("favourites API: %s", err.Error())
	}

	return faves, nil
}
