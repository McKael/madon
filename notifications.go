package gondole

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// GetNotifications returns the list of the user's notifications
func (g *Client) GetNotifications() ([]Notification, error) {
	var notifications []Notification

	req := g.prepareRequest("notifications")
	r, err := rest.API(req)
	if err != nil {
		return notifications, fmt.Errorf("notifications API query: %s", err.Error())
	}

	println(r.Body)
	err = json.Unmarshal([]byte(r.Body), &notifications)
	if err != nil {
		var res struct {
			Error string `json:"error"`
		}
		err2 := json.Unmarshal([]byte(r.Body), &res)
		if err2 == nil {
			return notifications, fmt.Errorf("%s", res.Error)
		}
		return notifications, fmt.Errorf("notifications API: %s", err.Error())
	}

	return notifications, nil
}

// GetNotification returns a notification
func (g *Client) GetNotification(id int) (*Notification, error) {
	var notification Notification

	req := g.prepareRequest("notifications/" + strconv.Itoa(id))
	r, err := rest.API(req)
	if err != nil {
		return &notification, fmt.Errorf("notification API query: %s", err.Error())
	}

	println(r.Body)
	err = json.Unmarshal([]byte(r.Body), &notification)
	if err != nil {
		var res struct {
			Error string `json:"error"`
		}
		err2 := json.Unmarshal([]byte(r.Body), &res)
		if err2 == nil {
			return &notification, fmt.Errorf("%s", res.Error)
		}
		return &notification, fmt.Errorf("notification API: %s", err.Error())
	}

	return &notification, nil
}
