package commands

import (
	"os"

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

		config := new(model.DeployServerConfig)
		err := helpers.LoadConfiguration(configName, config)
		if err != nil {
			return err
		}

		deployNameStyle := chalk.Cyan.NewStyle().WithTextStyle(chalk.Bold).Style
		fmt.Println(chalk.Blue.Color(fmt.Sprintf("Preparing %s ...", deployNameStyle(configName))))

		// ssh into the deployment server
		remoteName := config.RemoteName()
		branchName := config.BranchName()

		commands := []string{
			"rm -rf " + config.Path,
			"mkdir -p " + config.Path,
			fmt.Sprintf("cd %s && git clone -b %s -o %s --single-branch %s source", config.Path, branchName, remoteName, config.Repo),
		}

		_, err = helpers.ExecuteCmd(commands, config)
		helpers.FailOnError(err, "Failed to prepare remote server folder")

		fmt.Println(chalk.Green.Color("Succesfully prepared"))

		return nil
	},
}
