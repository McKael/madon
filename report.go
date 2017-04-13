package gondole

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/sendgrid/rest"
)

// GetReports returns the current user's reports
func (g *Client) GetReports() ([]Report, error) {
	req := g.prepareRequest("reports")
	r, err := rest.API(req)
	if err != nil {
		return nil, fmt.Errorf("reports: %s", err.Error())
	}

	// Check for error reply
	var errorResult Error
	if err := json.Unmarshal([]byte(r.Body), &errorResult); err == nil {
		// The empty object is not an error
		if errorResult.Text != "" {
			return nil, fmt.Errorf("%s", errorResult.Text)
		}
	}

	// Not an error reply; let's unmarshal the data
	var reports []Report
	err = json.Unmarshal([]byte(r.Body), &reports)
	if err != nil {
		return nil, fmt.Errorf("reports API: %s", err.Error())
	}
	return reports, nil
}

// ReportUser reports the user account
// NOTE: Currently only the first statusID is sent.
func (g *Client) ReportUser(accountID int, statusIDs []int, comment string) (*Report, error) {
	if accountID < 1 || comment == "" || len(statusIDs) < 1 {
		return nil, ErrInvalidParameter
	}

	req := g.prepareRequest("reports")
	req.Method = rest.Post
	req.QueryParams["account_id"] = strconv.Itoa(accountID)
	// XXX Sending only the first one since I'm not sure arrays work
	req.QueryParams["status_ids"] = strconv.Itoa(statusIDs[0])
	req.QueryParams["comment"] = comment

	r, err := rest.API(req)
	if err != nil {
		return nil, fmt.Errorf("reports: %s", err.Error())
	}

	// Check for error reply
	var errorResult Error
	if err := json.Unmarshal([]byte(r.Body), &errorResult); err == nil {
		// The empty object is not an error
		if errorResult.Text != "" {
			return nil, fmt.Errorf("%s", errorResult.Text)
		}
	}

	// Not an error reply; let's unmarshal the data
	var report Report
	err = json.Unmarshal([]byte(r.Body), &report)
	if err != nil {
		return nil, fmt.Errorf("reports API: %s", err.Error())
	}
	return &report, nil
}
