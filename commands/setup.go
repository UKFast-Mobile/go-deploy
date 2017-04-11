package commands

import (
	"github.com/urfave/cli"
)

// Setup provides setup command for the deployment configuration
var Setup = cli.Command{
	Name:  "setup",
	Usage: "Setup deployment configuration, adds to or creates configuration json file",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "overrides current configuration for a given name if already exists",
		},
	},
}
