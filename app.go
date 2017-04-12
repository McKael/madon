package gondole

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"

	"github.com/sendgrid/rest"
)

type registerApp struct {
	ID           int    `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NewApp registers a new instance
func NewApp(name string, scopes []string, redirectURI, instanceURL string) (g *Client, err error) {
	if instanceURL == "" {
		instanceURL = defaultInstanceURL
	}

	if !strings.Contains(instanceURL, "://") {
		instanceURL = "https://" + instanceURL
	}

	apiPath := instanceURL + defaultAPIPath

	if _, err := url.ParseRequestURI(apiPath); err != nil {
		return nil, err
	}

	g = &Client{
		Name:        name,
		APIBase:     apiPath,
		InstanceURL: instanceURL,
	}

	req := g.prepareRequest("apps")
	if redirectURI != "" {
		req.QueryParams["redirect_uris"] = redirectURI
	} else {
		req.QueryParams["redirect_uris"] = NoRedirect
	}
	req.QueryParams["client_name"] = name
	req.QueryParams["scopes"] = strings.Join(scopes, " ")
	req.Method = rest.Post

	r, err := rest.API(req)
	if err != nil {
		log.Fatalf("error can not register app: %v", err)
	}

	var resp registerApp

	err = json.Unmarshal([]byte(r.Body), &resp)
	if err != nil {
		log.Fatalf("error can not register app: %v", err)
	}

	g.ID = resp.ClientID
	g.Secret = resp.ClientSecret

	return
}
