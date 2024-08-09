package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

const defaultConfigPath = "~/.kube/config"

func readConfig(path string) (map[string]interface{}, error) {
	expandedPath := expandPath(path)

	// Check if the file exists
	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		// If it doesn't exist, create an empty config file
		fmt.Printf("Config file not found at %s. Creating an empty config...\n", expandedPath)
		if err := createConfigFromFile(expandedPath, ""); err != nil {
			return nil, fmt.Errorf("failed to create config file: %v", err)
		}
	}

	data, err := os.ReadFile(expandedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config map[string]interface{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return config, nil
}

func createConfigFromFile(configPath, inputPath string) error {
	// If inputPath is empty, create an empty config file
	if inputPath == "" {
		emptyConfig := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Config",
			"clusters":   []interface{}{},
			"contexts":   []interface{}{},
			"users":      []interface{}{},
		}
		return writeConfig(configPath, emptyConfig)
	}

	// Read the input file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the input file contents to the new config file
	if err := os.WriteFile(configPath, inputData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Created new config file at %s\n", configPath)
	return nil
}

func writeConfig(path string, config map[string]interface{}) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(path, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func backupConfig(path string) (string, error) {
	backupPath := path + "." + time.Now().Format("20060102150405") + ".bak"
	err := copyFile(path, backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %v", err)
	}
	return backupPath, nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0600)
}

func getClusterNames(configPath string) ([]string, error) {
	config, err := readConfig(configPath)
	if err != nil {
		return nil, err
	}

	var clusterNames []string
	if clusters, ok := config["clusters"].([]interface{}); ok {
		for _, cluster := range clusters {
			if c, ok := cluster.(map[interface{}]interface{}); ok {
				if name, ok := c["name"].(string); ok {
					clusterNames = append(clusterNames, name)
				}
			}
		}
	}

	sort.Strings(clusterNames)
	return clusterNames, nil
}

func addClusterConfig(mainConfigPath, newConfigPath string) error {
	mainConfig, err := readConfig(mainConfigPath)
	if err != nil {
		return err
	}

	newConfig, err := readConfig(newConfigPath)
	if err != nil {
		return err
	}

	backupPath, err := backupConfig(mainConfigPath)
	if err != nil {
		return err
	}

	// Merge new config into main config
	for _, key := range []string{"clusters", "contexts", "users"} {
		if newItems, ok := newConfig[key].([]interface{}); ok {
			if mainItems, ok := mainConfig[key].([]interface{}); ok {
				mainConfig[key] = append(mainItems, newItems...)
			} else {
				mainConfig[key] = newItems
			}
		}
	}

	err = writeConfig(mainConfigPath, mainConfig)
	if err != nil {
		return err
	}

	fmt.Printf("New cluster config added. Backup created at %s\n", backupPath)
	return nil
}

func removeClusterConfig(mainConfigPath, clusterName string) error {
	config, err := readConfig(mainConfigPath)
	if err != nil {
		return err
	}

	backupPath, err := backupConfig(mainConfigPath)
	if err != nil {
		return err
	}

	removed := false
	for _, key := range []string{"clusters", "contexts", "users"} {
		if items, ok := config[key].([]interface{}); ok {
			newItems := make([]interface{}, 0, len(items))
			for _, item := range items {
				if m, ok := item.(map[interface{}]interface{}); ok {
					if m["name"] != clusterName {
						newItems = append(newItems, item)
					} else {
						removed = true
					}
				}
			}
			config[key] = newItems
		}
	}

	if !removed {
		return fmt.Errorf("cluster %s not found", clusterName)
	}

	err = writeConfig(mainConfigPath, config)
	if err != nil {
		return err
	}

	fmt.Printf("Cluster %s removed. Backup created at %s\n", clusterName, backupPath)
	return nil
}

func rollbackConfig(mainConfigPath string) error {
	dir := filepath.Dir(mainConfigPath)
	base := filepath.Base(mainConfigPath)

	files, err := filepath.Glob(filepath.Join(dir, base+".*.bak"))
	if err != nil {
		return fmt.Errorf("failed to find backup files: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no backup files found")
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))
	latestBackup := files[0]

	err = copyFile(latestBackup, mainConfigPath)
	if err != nil {
		return fmt.Errorf("failed to restore from backup: %v", err)
	}

	fmt.Printf("Config rolled back to %s\n", latestBackup)
	return nil
}
