package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Version information set by build flags
var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	app := &cli.App{
		Name:                 "kubecnf",
		Usage:                "manage cluster configs in kubectl config",
		Version:              version,
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
			versionCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var versionCommand = &cli.Command{
	Name:  "version",
	Usage: "show version information",
	Action: func(c *cli.Context) error {
		fmt.Printf("kubecnf version %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		return nil
	},
}
