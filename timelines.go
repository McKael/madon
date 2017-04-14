package gondole

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sendgrid/rest"
)

// GetTimelines returns a timeline (a list of statuses
// timeline can be "home", "public", or a hashtag (use ":hashtag")
// For the public timelines, you can set 'local' to true to get only the
// local instance.
func (g *Client) GetTimelines(timeline string, local bool) ([]Status, error) {
	var endPoint string
	var tl []Status

	switch {
	case timeline == "home", timeline == "public":
		endPoint = "timelines/" + timeline
	case strings.HasPrefix(timeline, ":"):
		hashtag := timeline[1:]
		if hashtag == "" {
			return tl, fmt.Errorf("timelines API: empty hashtag")
		}
		endPoint = "timelines/tag/" + hashtag
	default:
		return tl, fmt.Errorf("GetTimelines: bad timelines argument")
	}

	req := g.prepareRequest(endPoint)

	if timeline == "public" && local {
		req.QueryParams["local"] = "true"
	}

	r, err := rest.API(req)
	if err != nil {
		return tl, fmt.Errorf("timelines API query: %s", err.Error())
	}

	err = json.Unmarshal([]byte(r.Body), &tl)
	if err != nil {
		var errorRes Error
		err2 := json.Unmarshal([]byte(r.Body), &errorRes)
		if err2 == nil {
			return tl, fmt.Errorf("%s", errorRes.Text)
		}
		return tl, fmt.Errorf("timelines API: %s", err.Error())
	}

	return tl, nil
}
