package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:                 "kubecnf",
		Usage:                "manage cluster configs in kubectl config",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "specify the main config file",
				Value:   defaultConfigPath,
			},
		},
		Commands: []*cli.Command{
			addCommand,
			removeCommand,
			listCommand,
			rollbackCommand,
			completionCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
