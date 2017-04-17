/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package gondole

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sendgrid/rest"
)

// StreamEvent contains a single event from the streaming API
type StreamEvent struct {
	Event string      // Name of the event (error, update, notification or delete)
	Data  interface{} // Status, Notification or status ID
	Error error       // Error message from the StreamListener
}

// openStream opens a stream URL and returns an http.Response
// Note that the caller should close the connection when it's done reading
// the stream.
// The stream name can be "user", "public" or "hashtag".
// For "hashtag", the hashTag argument cannot be empty.
func (g *Client) openStream(streamName, hashTag string) (*http.Response, error) {
	params := make(apiCallParams)

	switch streamName {
	case "user", "public":
	case "hashtag":
		if hashTag == "" {
			return nil, ErrInvalidParameter
		}
		params["tag"] = hashTag
	default:
		return nil, ErrInvalidParameter
	}

	req, err := g.prepareRequest("streaming/"+streamName, rest.Get, params)
	if err != nil {
		return nil, fmt.Errorf("cannot build stream request: %s", err.Error())
	}

	reqObj, err := rest.BuildRequestObject(req)
	if err != nil {
		return nil, fmt.Errorf("cannot build stream request: %s", err.Error())
	}

	resp, err := rest.MakeRequest(reqObj)
	if err != nil {
		return nil, fmt.Errorf("cannot open stream: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, errors.New(resp.Status)
	}
	return resp, nil
}

// readStream reads from the http.Response and sends events to the events channel
// It stops when the connection is closed or when the stopCh channel is closed.
// The foroutine will close the doneCh channel when it terminates.
func (g *Client) readStream(events chan<- StreamEvent, stopCh <-chan bool, doneCh chan<- bool, r *http.Response) {
	defer r.Body.Close()

	reader := bufio.NewReader(r.Body)

	var line, eventName string
	for {
		select {
		case <-stopCh:
			close(doneCh)
			return
		default:
		}

		lineBytes, partial, err := reader.ReadLine()
		if err != nil {
			var e error
			if err == io.EOF {
				e = fmt.Errorf("connection closed: %s", err.Error())
			} else {
				e = fmt.Errorf("read error: %s", err.Error())
			}
			log.Printf("Stream Reader: %s", e.Error())
			events <- StreamEvent{Event: "error", Error: e}
			close(doneCh)
			return
		}

		if partial {
			e := fmt.Errorf("received incomplete line; not supported yet")
			log.Printf("Stream Reader: %s", e.Error())
			events <- StreamEvent{Event: "error", Error: e}
			time.Sleep(5 * time.Second)
			continue // Skip this
		}

		line = string(bytes.TrimSpace(lineBytes))

		if line == "" {
			continue // Skip empty line
		}
		if strings.HasPrefix(line, ":") {
			continue // Skip comment
		}

		if strings.HasPrefix(line, "event: ") {
			eventName = line[7:]
			continue
		}

		if !strings.HasPrefix(line, "data: ") {
			// XXX Needs improvement
			e := fmt.Errorf("received unhandled event line '%s'", strings.Split(line, ":")[0])
			log.Printf("Stream Reader: %s", e.Error())
			events <- StreamEvent{Event: "error", Error: e}
			continue
		}

		// This is a data line
		data := []byte(line[6:])

		var obj interface{}

		// Decode API object
		switch eventName {
		case "update":
			var s Status
			if err := json.Unmarshal(data, &s); err != nil {
				e := fmt.Errorf("could not unmarshal data: %s", err.Error())
				log.Printf("Stream Reader: %s", e.Error())
				events <- StreamEvent{Event: "error", Error: e}
				continue
			}
			obj = s
		case "notification":
			var notif Notification
			if err := json.Unmarshal(data, &notif); err != nil {
				e := fmt.Errorf("could not unmarshal data: %s", err.Error())
				log.Printf("Stream Reader: %s", e.Error())
				events <- StreamEvent{Event: "error", Error: e}
				continue
			}
			obj = notif
		case "delete":
			var statusID int
			if err := json.Unmarshal(data, &statusID); err != nil {
				e := fmt.Errorf("could not unmarshal data: %s", err.Error())
				log.Printf("Stream Reader: %s", e.Error())
				events <- StreamEvent{Event: "error", Error: e}
				continue
			}
			obj = statusID
		case "":
			fallthrough
		default:
			e := fmt.Errorf("unhandled event '%s'", eventName)
			log.Printf("Stream Reader: %s", e.Error())
			events <- StreamEvent{Event: "error", Error: e}
			continue
		}

		// Send event to the channel
		events <- StreamEvent{Event: eventName, Data: obj}
	}
}

// StreamListener listens to a stream from the Mastodon server
// The stream 'name' can be "user", "public" or "hashtag".
// For 'hashtag', the hashTag argument cannot be empty.
// The events are sent to the events channel (the errors as well).
// The streaming is terminated if the 'stopCh' channel is closed.
// The 'doneCh' channel is closed if the connection is closed by the server.
// Please note that this method launches a goroutine to listen to the events.
func (g *Client) StreamListener(name, hashTag string, events chan<- StreamEvent, stopCh <-chan bool, doneCh chan<- bool) error {
	if g == nil {
		return fmt.Errorf("use of uninitialized gondole client")
	}

	resp, err := g.openStream(name, hashTag)
	if err != nil {
		return err
	}
	go g.readStream(events, stopCh, doneCh, resp)
	return nil
}
