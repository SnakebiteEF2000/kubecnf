package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var addCommand = &cli.Command{
	Name:      "add",
	Usage:     "add a new cluster config to the main config (from file or piped stdin)",
	ArgsUsage: "[input_file]",
	Action: func(c *cli.Context) error {
		mainConfigPath := expandPath(c.String("config"))

		// Check if data is being piped via stdin
		stdinStat, err := os.Stdin.Stat()
		if err != nil {
			return fmt.Errorf("failed to check stdin: %v", err)
		}

		// Check if stdin is not a character device (i.e., it's piped or redirected)
		isPiped := (stdinStat.Mode() & os.ModeCharDevice) == 0
		hasFileArg := c.NArg() >= 1

		// Validate input method
		if !isPiped && !hasFileArg {
			return fmt.Errorf("input file is required or pipe kubeconfig data to stdin")
		}

		if isPiped && hasFileArg {
			return fmt.Errorf("cannot use both file argument and piped input at the same time")
		}

		// Handle main config file creation if it doesn't exist
		if _, err := os.Stat(mainConfigPath); os.IsNotExist(err) {
			if isPiped {
				fmt.Printf("Main config file not found at %s. Creating it from piped input...\n", mainConfigPath)
				if err := createConfigFromStdin(mainConfigPath); err != nil {
					return fmt.Errorf("failed to create main config file from piped input: %v", err)
				}
				fmt.Printf("Main config file created at %s\n", mainConfigPath)
				return nil
			} else {
				fmt.Printf("Main config file not found at %s. Creating it from the input file...\n", mainConfigPath)
				newConfigPath := c.Args().First()
				if err := createConfigFromFile(mainConfigPath, newConfigPath); err != nil {
					return fmt.Errorf("failed to create main config file: %v", err)
				}
				fmt.Printf("Main config file created at %s\n", mainConfigPath)
				return nil
			}
		}

		// Add cluster config from piped input or file
		if isPiped {
			return addClusterConfigFromStdin(mainConfigPath)
		} else {
			newConfigPath := c.Args().First()
			return addClusterConfig(mainConfigPath, newConfigPath)
		}
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

var completionCommand = &cli.Command{
	Name:  "completion",
	Usage: "output shell completion code",
	Action: func(c *cli.Context) error {
		fmt.Print(bashCompletionScript)
		return nil
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
