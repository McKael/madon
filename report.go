package gondole

import (
	"strconv"

	"github.com/sendgrid/rest"
)

// GetReports returns the current user's reports
func (g *Client) GetReports() ([]Report, error) {
	var reports []Report
	if err := g.apiCall("reports", rest.Get, nil, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}

// ReportUser reports the user account
// NOTE: Currently only the first statusID is sent.
func (g *Client) ReportUser(accountID int, statusIDs []int, comment string) (*Report, error) {
	if accountID < 1 || comment == "" || len(statusIDs) < 1 {
		return nil, ErrInvalidParameter
	}

	params := make(apiCallParams)
	params["account_id"] = strconv.Itoa(accountID)
	// XXX Sending only the first one since I'm not sure arrays work
	params["status_ids"] = strconv.Itoa(statusIDs[0])
	params["comment"] = comment

	var report Report
	if err := g.apiCall("reports", rest.Post, params, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
