package main

import (
	"os"

	"github.com/UKFast-Mobile/go-deploy/commands"
	"github.com/UKFast-Mobile/go-deploy/helpers"
	"github.com/UKFast-Mobile/go-deploy/model"
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
		commands.Prepare,
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

		// deployNameStyle := chalk.Cyan.NewStyle().WithTextStyle(chalk.Bold).Style

		// fmt.Println(chalk.Blue.Color(fmt.Sprintf("Deploying to %s ...", deployNameStyle(configName))))

		// // ssh into the deployment server
		// sshConfig := &ssh.ClientConfig{
		// 	User: config.Username,
		// 	Auth: []ssh.AuthMethod{
		// 		helpers.PublicKeyFile(config.PrivateKey),
		// 	},
		// }

		// connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", config.Host, config.Port), sshConfig)
		// helpers.FailOnError(err, "Failed to dial deployment server")
		// defer connection.Close()

		// session, err := connection.NewSession()
		// helpers.FailOnError(err, "Failed to create a session with the deployment server")
		// defer session.Close()

		// fmt.Println(chalk.Blue.Color(fmt.Sprintf("On %s", deployNameStyle(configName))))

		// err = session.Run(fmt.Sprintf("mkdir -p %s", config.Path))
		// helpers.FailOnError(err, "Failed to create deployment folder")

		// err = session.Run(fmt.Sprintf("git clone %s"))

		return nil
	}

	app.Run(os.Args)
}
