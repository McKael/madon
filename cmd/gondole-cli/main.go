package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/McKael/gondole"
)

var (
	fVerbose             bool
	fInstance            string
	fAuthMethod          string
	fUsername, fPassword string
	fScopes              string

	instance *gondole.Client
	cnf      *Server

	// Default scopes
	ourScopes = []string{
		"read",
		"write",
		"follow",
	}

	defaultInstanceURL = "https://mastodon.social"

	authMethods = map[string]bool{
		"basic":  true,
		"oauth2": true,
	}
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

	// Can be "oauth2", "basic"
	Auth string

	// If not using OAuth2
	User     string
	Password string
}

func setupEnvironment(c *cli.Context) (err error) {
	var config Config
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

	if fAuthMethod != "" && authMethods[fAuthMethod] {

	}

	// Load configuration, will register if none is found
	cnf, err = LoadConfig(instanceName)
	if err != nil {
		// Nothing exist yet
		/*
			defName := Config{
				Default:  instanceName,
				Auth:     "basic",
				User:     "",
				Password: "",
			}
		*/

		err = config.Write()
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
	if fAuthMethod == "basic" {
		err = instance.LoginBasic(fUsername, fPassword)
	}

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
		cli.StringFlag{
			Name:        "auth,A",
			Usage:       "authentication mode",
			Destination: &fAuthMethod,
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
		cli.StringFlag{
			Name:        "username,login",
			Usage:       "user name",
			Destination: &fUsername,
		},
		cli.StringFlag{
			Name:        "password",
			Usage:       "user password",
			Destination: &fPassword,
		},
		cli.BoolFlag{
			Name:        "verbose,v",
			Usage:       "verbose mode",
			Destination: &fVerbose,
		},
	}
	app.Run(os.Args)
}
