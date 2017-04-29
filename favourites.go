/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"github.com/sendgrid/rest"
)

// GetFavourites returns the list of the user's favourites
func (mc *Client) GetFavourites(lopt *LimitParams) ([]Status, error) {
	var faves []Status
	var links apiLinks
	err := mc.apiCall("favourites", rest.Get, nil, lopt, &links, &faves)
	if err != nil {
		return nil, err
	}
	if lopt != nil { // Fetch more pages to reach our limit
		var faveSlice []Status
		for lopt.Limit > len(faves) && links.next != nil {
			newlopt := links.next
			links = apiLinks{}
			if err := mc.apiCall("favourites", rest.Get, nil, newlopt, &links, &faveSlice); err != nil {
				return nil, err
			}
			faves = append(faves, faveSlice...)
			faveSlice = faveSlice[:0] // Clear struct
		}
	}
	return faves, nil
}
