package main

import (
    "log"
    "github.com/keltia/gondole"
    "github.com/urfave/cli"
    "os"
)

var (
    fVerbose bool
)

func Register(c *cli.Context) error {

    return nil
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
    }
    app.Run(os.Args)
}
