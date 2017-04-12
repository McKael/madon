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
		var res struct {
			Error string `json:"error"`
		}
		err2 := json.Unmarshal([]byte(r.Body), &res)
		if err2 == nil {
			return faves, fmt.Errorf("%s", res.Error)
		}
		return faves, fmt.Errorf("favourites API: %s", err.Error())
	}

	return faves, nil
}
