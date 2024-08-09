package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var addCommand = &cli.Command{
	Name:      "add",
	Usage:     "add a new cluster config to the main config",
	ArgsUsage: "<input_file>",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return fmt.Errorf("input file is required")
		}
		mainConfigPath := expandPath(c.String("config"))
		newConfigPath := c.Args().First()

		// Check if the main config file exists
		if _, err := os.Stat(mainConfigPath); os.IsNotExist(err) {
			// If it doesn't exist, create it from the input file
			fmt.Printf("Main config file not found at %s. Creating it from the input file...\n", mainConfigPath)
			if err := createConfigFromFile(mainConfigPath, newConfigPath); err != nil {
				return fmt.Errorf("failed to create main config file: %v", err)
			}
			fmt.Printf("Main config file created at %s\n", mainConfigPath)
			return nil
		}

		// If the main config file exists, proceed with adding the new config
		return addClusterConfig(mainConfigPath, newConfigPath)
	},
}

var removeCommand = &cli.Command{
	Name:      "remove",
	Usage:     "remove a cluster config from the main config",
	ArgsUsage: "<cluster_name>",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return fmt.Errorf("cluster name is required")
		}
		mainConfigPath := expandPath(c.String("config"))
		clusterName := c.Args().First()
		return removeClusterConfig(mainConfigPath, clusterName)
	},
	BashComplete: func(c *cli.Context) {
		mainConfigPath := expandPath(c.String("config"))
		clusters, err := getClusterNames(mainConfigPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting cluster names: %v\n", err)
			return
		}
		for _, cluster := range clusters {
			fmt.Println(cluster)
		}
	},
}

var listCommand = &cli.Command{
	Name:  "list",
	Usage: "list all cluster configurations",
	Action: func(c *cli.Context) error {
		mainConfigPath := expandPath(c.String("config"))
		return listClusterConfigs(mainConfigPath)
	},
}

var rollbackCommand = &cli.Command{
	Name:  "rollback",
	Usage: "rollback to the previous config",
	Action: func(c *cli.Context) error {
		mainConfigPath := expandPath(c.String("config"))
		return rollbackConfig(mainConfigPath)
	},
}

func listClusterConfigs(mainConfigPath string) error {
	clusters, err := getClusterNames(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get cluster names: %v", err)
	}

	if len(clusters) == 0 {
		fmt.Println("No cluster configurations found.")
		return nil
	}

	fmt.Println("Existing cluster configurations:")
	for _, cluster := range clusters {
		fmt.Printf("- %s\n", cluster)
	}

	return nil
}
