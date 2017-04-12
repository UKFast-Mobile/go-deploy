package main

import (
	"os"

	"github.com/UKFast-Mobile/go-deploy/commands"
	"gopkg.in/urfave/cli.v1"
)

// Debug - show logs
var Debug bool

func main() {
	app := cli.NewApp()

	app.Name = "Go Deploy"
	app.Usage = "Deploy your docker applications"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aleksandr Kelbas",
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Show debug logs (verbose mode)",
			EnvVar:      "DEBUG",
			Destination: &Debug,
		},
	}

	app.Commands = []cli.Command{
		commands.Setup,
	}
	app.Run(os.Args)
}
