package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"

	"log"

	"github.com/UKFast-Mobile/go-deploy/commands"
	"github.com/UKFast-Mobile/go-deploy/helpers"
	"github.com/UKFast-Mobile/go-deploy/model"
	"github.com/ttacon/chalk"
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

	app.Action = func(c *cli.Context) error {

		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		configName := c.Args()[0]

		var config model.DeployServerConfig
		err := helpers.LoadConfiguration(configName, &config)
		if err != nil {
			return err
		}

		deployNameStyle := chalk.Cyan.NewStyle().WithTextStyle(chalk.Bold).Style

		fmt.Println(chalk.Blue.Color(fmt.Sprintf("Deploying to %s ...", deployNameStyle(configName))))

		// ssh into the deployment server
		sshConfig := &ssh.ClientConfig{
			User: config.Username,
			Auth: []ssh.AuthMethod{
				helpers.PublicKeyFile(config.PrivateKey),
			},
		}

		connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), sshConfig)
		helpers.FailOnError(err, "Failed to dial deployment server")
		defer connection.Close()

		session, err := connection.NewSession()
		helpers.FailOnError(err, "Failed to create a session with the deployment server")
		defer session.Close()

		data, err := session.Output("pwd")
		helpers.FailOnError(err, "Failed to run pwd")

		result := string(data)

		log.Printf("%s: %s", config.Host, result)

		return nil
	}

	app.Run(os.Args)
}
