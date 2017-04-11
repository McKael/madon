package gondole

import (
	"encoding/json"
	"github.com/sendgrid/rest"
	"log"
	"strings"
)

var (
	ourScopes = []string{
		"read",
		"write",
		"follow",
	}
)

type registerApp struct {
	ID           string `json:"id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func registerApplication(name string, scopes []string, redirectURI string) (g *Gondole, err error) {
	g = &Gondole{
		Name: name,
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

	server := &Server{
		ID:          g.ID,
		Name:        name,
		BearerToken: g.Secret,
	}
	err = server.WriteToken(name)
	if err != nil {
		log.Fatalf("error: can not write token for %s", name)
	}
	return
}

// NewApp registers a new instance
func NewApp(name string, scopes []string, redirectURI string) (g *Gondole, err error) {
	// Load configuration, will register if none is found
	cnf, err := LoadConfig(name)
	if err != nil {
		// Nothing exist yet
		cnf := Config{
			Default: name,
		}
		err = cnf.Write()
		if err != nil {
			log.Fatalf("error: can not write config for %s", name)
		}

		// Now register this through OAuth
		if scopes == nil {
			scopes = ourScopes
		}

		g, err = registerApplication(name, scopes, redirectURI)

	} else {
		g = &Gondole{
			Name:   cnf.Name,
			ID:     cnf.ID,
			Secret: cnf.BearerToken,
		}
	}

	return
}
