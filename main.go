package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "Go Deploy"
	app.Usage = "Deploy your docker applications"
	app.Version = "0.0.1"

	app.Run(os.Args)
}
