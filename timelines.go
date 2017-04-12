package gondole

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sendgrid/rest"
)

// GetTimelines returns a timeline (a list of statuses
// timeline can be "home", "public", or a hashtag (":hashtag")
func (g *Client) GetTimelines(timeline string) ([]Status, error) {
	var endPoint string
	var tl []Status

	if timeline == "home" || timeline == "public" {
		endPoint = "timelines/" + timeline
	} else if strings.HasPrefix(timeline, ":") {
		hashtag := timeline[1:]
		if hashtag == "" {
			return tl, fmt.Errorf("timelines API: empty hashtag")
		}
		endPoint = "timelines/tag/" + hashtag
	} else {
		return tl, fmt.Errorf("GetTimelines: bad timelines argument")
	}

	req := g.prepareRequest(endPoint)
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
