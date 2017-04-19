package commands

import (
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"fmt"

	"github.com/UKFast-Mobile/go-deploy/helpers"
	"github.com/UKFast-Mobile/go-deploy/model"
	"github.com/ttacon/chalk"
	"gopkg.in/urfave/cli.v1"
)

// Prepare prepares remote server folder (clones repo)
var Prepare = cli.Command{
	Name:  "prepare",
	Usage: "Prepare remote server from provided configuration file",
	Action: func(c *cli.Context) error {
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
		fmt.Println(chalk.Blue.Color(fmt.Sprintf("Setting up %s ...", deployNameStyle(configName))))

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

		fmt.Println(chalk.Blue.Color(fmt.Sprintf("On %s", deployNameStyle(configName))))

		err = session.Run(fmt.Sprintf("rm -rf %s", config.Path))
		helpers.FailOnError(err, "Failed to remove old directory")

		err = session.Run(fmt.Sprintf("mkdir -p %s", config.Path))
		helpers.FailOnError(err, "Failed to create deployment folder")

		err = session.Run(fmt.Sprintf("cd %s", config.Path))
		helpers.FailOnError(err, "Failed to navigate to deployment folder")

		refs := strings.Split(config.Ref, "/")
		remoteName := refs[0]
		branchName := refs[1]

		err = session.Run(fmt.Sprintf("git clone -b %s -o %s %s source", branchName, remoteName, config.Repo))
		helpers.FailOnError(err, "Failed to clone remote repository")

		return nil
	},
}
