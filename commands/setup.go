package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/UKFast-Mobile/go-deploy/helpers"
	"github.com/UKFast-Mobile/go-deploy/model"

	"reflect"

	"github.com/ttacon/chalk"
	"gopkg.in/urfave/cli.v1"
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
	Action: func(c *cli.Context) {
		if c.NArg() == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
			// return errors.New("Invalid number of arguments")
		}

		force := c.Bool("force")
		if force {
			fmt.Println(chalk.Red, "! running in force mode will override current configuration (if exists)")
		}

		configName := c.Args()[0]
		fmt.Println(chalk.Magenta, "Setting up ", chalk.Bold.TextStyle(configName), " deployment configuration")

		// Load the configuration file
		file, err := helpers.OpenConfigFile()
		if file[configName] != nil && !force {
			fmt.Println(chalk.Red, "Failed: configuration already exists!")
			os.Exit(1)
		}

		config := model.DeployServerConfig{}
		t := reflect.TypeOf(config)

		scanner := bufio.NewScanner(os.Stdin)
		log.Printf("Num fields: %d", t.NumField())

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("cli_q")
			if tag != "" {
				fmt.Printf(tag)
				scanner.Scan()
				reflect.ValueOf(&config).Elem().Field(i).SetString(scanner.Text())
			}
		}

		file[configName] = config
		err = helpers.WriteToFile(file)
		helpers.FailOnError(err, "Failed to write to a config file")

		os.Exit(0)
	},
}
