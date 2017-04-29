/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"fmt"
	"strings"

	"github.com/sendgrid/rest"
)

// GetTimelines returns a timeline (a list of statuses
// timeline can be "home", "public", or a hashtag (use ":hashtag" or "#hashtag")
// For the public timelines, you can set 'local' to true to get only the
// local instance.
// If lopt.All is true, several requests will be made until the API server
// has nothing to return.
// If lopt.Limit is set (and not All), several queries can be made until the
// limit is reached.
func (mc *Client) GetTimelines(timeline string, local bool, lopt *LimitParams) ([]Status, error) {
	var endPoint string

	switch {
	case timeline == "home", timeline == "public":
		endPoint = "timelines/" + timeline
	case strings.HasPrefix(timeline, ":"), strings.HasPrefix(timeline, "#"):
		hashtag := timeline[1:]
		if hashtag == "" {
			return nil, fmt.Errorf("timelines API: empty hashtag")
		}
		endPoint = "timelines/tag/" + hashtag
	default:
		return nil, fmt.Errorf("GetTimelines: bad timelines argument")
	}

	params := make(apiCallParams)
	if timeline == "public" && local {
		params["local"] = "true"
	}

	var tl []Status
	var links apiLinks
	if err := mc.apiCall(endPoint, rest.Get, params, lopt, &links, &tl); err != nil {
		return nil, err
	}
	if lopt != nil { // Fetch more pages to reach our limit
		var statusSlice []Status
		for (lopt.All || lopt.Limit > len(tl)) && links.next != nil {
			newlopt := links.next
			links = apiLinks{}
			if err := mc.apiCall(endPoint, rest.Get, params, newlopt, &links, &statusSlice); err != nil {
				return nil, err
			}
			tl = append(tl, statusSlice...)
			statusSlice = statusSlice[:0] // Clear struct
		}
	}
	return tl, nil
}
