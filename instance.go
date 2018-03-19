/*
Copyright 2017-2018 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"github.com/sendgrid/rest"
)

// GetCurrentInstance returns current instance information
func (mc *Client) GetCurrentInstance() (*Instance, error) {
	var i Instance
	if err := mc.apiCall("instance", rest.Get, nil, nil, nil, &i); err != nil {
		return nil, err
	}
	return &i, nil
}

// GetInstancePeers returns current instance peers
func (mc *Client) GetInstancePeers() ([]string, error) {
	var peers []string
	if err := mc.apiCall("instance/peers", rest.Get, nil, nil, nil, &peers); err != nil {
		return nil, err
	}
	return peers, nil
}

// GetInstanceActivity returns current instance activity
func (mc *Client) GetInstanceActivity() (interface{}, error) {
	var activity interface{}
	if err := mc.apiCall("instance/activity", rest.Get, nil, nil, nil, &activity); err != nil {
		return nil, err
	}
	return activity, nil
}
