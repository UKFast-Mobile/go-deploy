package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/UKFast-Mobile/go-deploy/model"

	"github.com/ttacon/chalk"

	"reflect"

	"io/ioutil"

	"encoding/json"

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
		
		configName := c.Args()[0]
		fmt.Printf("Did get config name: %s\n", configName)




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

		// var file *os.File
		// fileInfo, err := os.Stat("./go-deploy.json")

		file, err := os.OpenFile("./go-deploy.json", os.O_RDWR|os.O_CREATE, 0666)
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		var fileJSON map[string]interface{}
		_ = json.Unmarshal(data, &fileJSON)

		if fileJSON[configName] != nil && !force {
			// config already exists
			log.Fatal("Config with given  name already exists!")
		} else {

			if fileJSON == nil {
				fileJSON = map[string]interface{}{}
			}

			fileJSON[configName] = config
			backToJSON, err := json.MarshalIndent(fileJSON, "", "  ")
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			err = ioutil.WriteFile("./go-deploy.json", backToJSON, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}

		os.Exit(0)
	},
}
