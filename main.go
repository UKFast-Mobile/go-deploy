package main

import (
	"fmt"
	"os"

	"github.com/UKFast-Mobile/go-deploy/commands"
	"github.com/UKFast-Mobile/go-deploy/helpers"
	"github.com/UKFast-Mobile/go-deploy/model"
	"github.com/ttacon/chalk"
	"gopkg.in/urfave/cli.v1"
)

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
			Name:        "debug, d",
			Usage:       "Show debug logs (verbose mode)",
			EnvVar:      "DEBUG",
			Destination: &helpers.Debug,
		},
	}

	app.Commands = []cli.Command{
		commands.Setup,
		commands.Prepare,
	}

	app.Action = func(c *cli.Context) error {

		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		configName := c.Args()[0]

		config := new(model.DeployServerConfig)
		err := helpers.LoadConfiguration(configName, config)
		if err != nil {
			return err
		}

		// Set env vars if any
		helpers.SetEnvVars(config)

		deployNameStyle := chalk.Cyan.NewStyle().WithTextStyle(chalk.Bold).Style

		fmt.Println(chalk.Blue.Color(fmt.Sprintf("Deploying to %s ...", deployNameStyle(configName))))

		// ssh into the deployment server
		commands := []string{
			fmt.Sprintf("cd %s/source && git checkout -q %s && git pull -q %s %s", config.Path, config.BranchName(), config.RemoteName(), config.BranchName()),
			fmt.Sprintf("cd %s/source && %s", config.Path, config.Cmd),
		}

		output, err := helpers.ExecuteCmd(commands, config)
		helpers.LogDebug(output)
		helpers.FailOnError(err, "Failed to deploy to server")

		fmt.Println(chalk.Green.Color("Deployed successfully"))

		return nil
	}

	app.Run(os.Args)
}
