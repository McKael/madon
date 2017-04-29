/*
Copyright 2017 Mikael Berthe

Licensed under the MIT license.  Please see the LICENSE file is this directory.
*/

package madon

import (
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// GetReports returns the current user's reports
// (I don't know if the limit options are used by the API server.)
func (mc *Client) GetReports(lopt *LimitParams) ([]Report, error) {
	var reports []Report
	if err := mc.apiCall("reports", rest.Get, nil, lopt, nil, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

// ReportUser reports the user account
func (mc *Client) ReportUser(accountID int, statusIDs []int, comment string) (*Report, error) {
	if accountID < 1 || comment == "" || len(statusIDs) < 1 {
		return nil, ErrInvalidParameter
	}

	params := make(apiCallParams)
	params["account_id"] = strconv.Itoa(accountID)
	params["comment"] = comment
	for i, id := range statusIDs {
		if id < 1 {
			return nil, ErrInvalidID
		}
		qID := fmt.Sprintf("status_ids[%d]", i+1)
		params[qID] = strconv.Itoa(id)
	}

	var report Report
	if err := mc.apiCall("reports", rest.Post, params, nil, nil, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
