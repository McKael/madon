package gondole

import (
	"encoding/json"
	"github.com/sendgrid/rest"
	"log"
	"strings"
)

var ()

type registerApp struct {
	ID           string `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// NewApp registers a new instance
func NewApp(name string, scopes []string, redirectURI, baseURL string) (g *Client, err error) {
	var endpoint string

	if baseURL != "" {
		endpoint = baseURL
	}

	g = &Client{
		Name: name,
		APIBase: endpoint,
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

	if err != nil {
		log.Fatalf("error: can not write token for %s", name)
	}

	g = &Client{
		Name:    name,
		ID:      resp.ClientID,
		Secret:  resp.ClientSecret,
		APIBase: endpoint,
	}

	return
}
