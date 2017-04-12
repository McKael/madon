package main

import (
	"github.com/keltia/gondole"
	"github.com/urfave/cli"
	"log"
	"os"
	"strings"
)

var (
	fVerbose  bool
	fInstance string
	fScopes   string

	instance *gondole.Client
	cnf      *Server

	// Default scopes
	ourScopes = []string{
		"read",
		"write",
		"follow",
	}

	defaultInstanceURL = "https://mastodon.social"
)

// Server holds our application details
type Server struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	BearerToken string `json:"bearer_token"`
	APIBase     string `json:"base_url"`
	InstanceURL string `json:"base_url"`
}

type Config struct {
	Default string
}

func setupEnvironment(c *cli.Context) (err error) {
	var scopes []string

	instanceURL := defaultInstanceURL
	if fInstance != "" {
		if strings.Contains(fInstance, "://") {
			instanceURL = fInstance
		} else {
			instanceURL = "https://" + fInstance
		}
	}

	instanceName := basename(instanceURL)

	// Load configuration, will register if none is found
	cnf, err = LoadConfig(instanceName)
	if err != nil {
		// Nothing exist yet
		defName := Config{
			Default: instanceName,
		}
		err = defName.Write()
		if err != nil {
			log.Fatalf("error: can not write config for %s", instanceName)
		}

		// Now register this through OAuth
		if fScopes != "" {
			scopes = strings.Split(fScopes, " ")
		} else {
			scopes = ourScopes
		}

		instance, err = gondole.NewApp("gondole-cli", scopes, gondole.NoRedirect, instanceURL)

		server := &Server{
			ID:          instance.ID,
			Name:        instance.Name,
			BearerToken: instance.Secret,
			APIBase:     instance.APIBase,
			InstanceURL: instance.InstanceURL,
		}
		err = server.WriteToken(instanceName)
		if err != nil {
			log.Fatalf("error: can not write token for %s", instance.Name)
		}

		cnf := Config{
			Default: instance.Name,
		}

		err = cnf.Write()
		if err != nil {
			log.Fatalf("error: can not write config for %s", instance.Name)
		}

	}
	// Log in to the instance
	return err
}

func init() {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	cli.VersionPrinter = func(c *cli.Context) {
		log.Printf("API wrapper: %s Mastodon CLI: %s\n", c.App.Version, gondole.GondoleVersion)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "gondole"
	app.Usage = "Mastodon CLI interface"
	app.Author = "Ollivier Robert <roberto@keltia.net>"
	app.Version = gondole.GondoleVersion
	//app.HideVersion = true

	app.Before = setupEnvironment

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose,v",
			Usage:       "verbose mode",
			Destination: &fVerbose,
		},
		cli.StringFlag{
			Name:        "instance,I",
			Usage:       "use that instance",
			Destination: &fInstance,
		},
		cli.StringFlag{
			Name:        "scopes,S",
			Usage:       "use these scopes",
			Destination: &fScopes,
		},
	}
	app.Run(os.Args)
}
