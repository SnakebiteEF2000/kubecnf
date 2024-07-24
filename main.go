package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

const defaultConfigPath = "~/.kube/config"

func main() {
	app := &cli.App{
		Name:  "kubecnf",
		Usage: "add or remove a cluster config from kubectl config",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "specify the main config file",
				Value:   defaultConfigPath,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "add a new cluster config to the main config",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "specify the new cluster config file",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					mainConfigPath := expandPath(c.String("config"))
					newConfigPath := c.String("file")
					return addClusterConfig(mainConfigPath, newConfigPath)
				},
			},
			{
				Name:  "remove",
				Usage: "remove a cluster config from the main config",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("cluster name is required")
					}
					mainConfigPath := expandPath(c.String("config"))
					clusterName := c.Args().Get(0)
					return removeClusterConfig(mainConfigPath, clusterName)
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback to the previous config",
				Action: func(c *cli.Context) error {
					mainConfigPath := expandPath(c.String("config"))
					return rollbackConfig(mainConfigPath)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func expandPath(path string) string {
	if path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func backupConfig(mainConfigPath string) (string, error) {
	backupPath := mainConfigPath + ".bak"
	err := copyFile(mainConfigPath, backupPath)
	if err != nil {
		return "", err
	}
	return backupPath, nil
}

func rollbackConfig(mainConfigPath string) error {
	backupPath := mainConfigPath + ".bak"
	return copyFile(backupPath, mainConfigPath)
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dst, input, 0600)
	if err != nil {
		return err
	}
	return nil
}

func addClusterConfig(mainConfigPath, newConfigPath string) error {
	backupPath, err := backupConfig(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	mainData, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read main config: %v", err)
	}

	newData, err := os.ReadFile(newConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read new config: %v", err)
	}

	var mainConfig, newConfig map[string]interface{}
	err = yaml.Unmarshal(mainData, &mainConfig)
	if err != nil {
		return fmt.Errorf("failed to parse main config: %v", err)
	}

	err = yaml.Unmarshal(newData, &newConfig)
	if err != nil {
		return fmt.Errorf("failed to parse new config: %v", err)
	}

	mergeConfigs(mainConfig, newConfig)

	updatedData, err := yaml.Marshal(&mainConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}

	err = os.WriteFile(mainConfigPath, updatedData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write updated config: %v", err)
	}

	fmt.Printf("Backup created at %s\n", backupPath)
	fmt.Println("Cluster added successfully.")
	return nil
}

func removeClusterConfig(mainConfigPath, clusterName string) error {
	backupPath, err := backupConfig(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	mainData, err := os.ReadFile(mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read main config: %v", err)
	}

	var mainConfig map[string]interface{}
	err = yaml.Unmarshal(mainData, &mainConfig)
	if err != nil {
		return fmt.Errorf("failed to parse main config: %v", err)
	}

	removeCluster(mainConfig, clusterName)

	updatedData, err := yaml.Marshal(&mainConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}

	err = os.WriteFile(mainConfigPath, updatedData, 0600)
	if err != nil {
		return fmt.Errorf("failed to write updated config: %v", err)
	}

	fmt.Printf("Backup created at %s\n", backupPath)
	fmt.Println("Cluster removed successfully.")
	return nil
}

func mergeConfigs(mainConfig, newConfig map[string]interface{}) {
	mergeSection(mainConfig, newConfig, "clusters")
	mergeSection(mainConfig, newConfig, "contexts")
	mergeSection(mainConfig, newConfig, "users")
}

func mergeSection(mainConfig, newConfig map[string]interface{}, section string) {
	if newSection, ok := newConfig[section].([]interface{}); ok {
		if mainSection, ok := mainConfig[section].([]interface{}); ok {
			mainConfig[section] = append(mainSection, newSection...)
		} else {
			mainConfig[section] = newSection
		}
	}
}

func removeCluster(config map[string]interface{}, clusterName string) {
	// Remove from clusters
	if clusters, ok := config["clusters"].([]interface{}); ok {
		for i := len(clusters) - 1; i >= 0; i-- {
			cluster := clusters[i].(map[interface{}]interface{})
			if cluster["name"] == clusterName {
				config["clusters"] = append(clusters[:i], clusters[i+1:]...)
			}
		}
	}

	// Remove from contexts
	if contexts, ok := config["contexts"].([]interface{}); ok {
		for i := len(contexts) - 1; i >= 0; i-- {
			context := contexts[i].(map[interface{}]interface{})
			if context["name"] == clusterName {
				config["contexts"] = append(contexts[:i], contexts[i+1:]...)
			}
		}
	}

	// Remove from users
	if users, ok := config["users"].([]interface{}); ok {
		for i := len(users) - 1; i >= 0; i-- {
			user := users[i].(map[interface{}]interface{})
			if user["name"] == clusterName {
				config["users"] = append(users[:i], users[i+1:]...)
			}
		}
	}

	// Update current-context so kubectl wont crsh
	if currentContext, ok := config["current-context"].(string); ok && currentContext == clusterName {
		config["current-context"] = ""
	}
}
