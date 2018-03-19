/*
Copyright 2017-2018 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"strings"

	"github.com/pkg/errors"
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
			return nil, errors.New("timelines API: empty hashtag")
		}
		endPoint = "timelines/tag/" + hashtag
	default:
		return nil, errors.New("GetTimelines: bad timelines argument")
	}

	params := make(apiCallParams)
	if timeline == "public" && local {
		params["local"] = "true"
	}

	return mc.getMultipleStatuses(endPoint, params, lopt)
}
