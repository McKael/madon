package main

import (
	"github.com/keltia/gondole"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	fVerbose bool
	fBaseURL string
    instance *gondole.Gondole
	cnf      *gondole.Config
)

func Register(c *cli.Context) (err error) {
    instance, err = gondole.NewApp("gondole-cli", nil, gondole.NoRedirect, fBaseURL)
	return err
}

func init() {
	cli.VersionFlag = cli.BoolFlag{Name: "version, V"}

	cli.VersionPrinter = func(c *cli.Context) {
		log.Printf("API wrapper: %s Mastodon CLI: %s\n", c.App.Version, gondole.APIVersion)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "gondole"
	app.Usage = "Mastodon CLI interface"
	app.Author = "Ollivier Robert <roberto@keltia.net>"
	app.Version = gondole.APIVersion
	//app.HideVersion = true

	app.Before = Register

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "verbose,v",
			Usage:       "verbose mode",
			Destination: &fVerbose,
		},
		cli.StringFlag{
			Name:        "instance,i",
			Usage:       "use that instance",
			Destination: &fBaseURL,
		},
	}
	app.Run(os.Args)
}
