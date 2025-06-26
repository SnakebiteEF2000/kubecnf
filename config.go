package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigPath = "~/.kube/config"
	// File and directory permissions
	dirPerm  = 0700 // Owner read/write/execute only
	filePerm = 0600 // Owner read/write only
)

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

	data, err := os.ReadFile(expandedPath) //nolint:gosec // Config path is user controlled
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
	inputData, err := os.ReadFile(inputPath) //nolint:gosec // File path is user provided
	if err != nil {
		return fmt.Errorf("failed to read input file: %v", err)
	}

	// Validate the input file as a proper kubeconfig
	var testConfig map[string]interface{}
	if err := yaml.Unmarshal(inputData, &testConfig); err != nil {
		return fmt.Errorf("invalid YAML in input file: %v", err)
	}

	if err := validateKubeconfig(testConfig); err != nil {
		return fmt.Errorf("invalid kubeconfig format in input file: %v", err)
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the input file contents to the new config file
	if err := os.WriteFile(configPath, inputData, filePerm); err != nil {
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

	err = os.WriteFile(path, data, filePerm)
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
	input, err := os.ReadFile(src) //nolint:gosec // File path is controlled by application
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, filePerm)
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

	// Merge new config into main config with duplicate detection
	err = mergeConfigs(mainConfig, newConfig)
	if err != nil {
		return err
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

	removed := removeClusterFromConfig(config, clusterName)

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

// removeClusterFromConfig removes cluster entries from all sections and returns true if any were found
func removeClusterFromConfig(config map[string]interface{}, clusterName string) bool {
	removed := false
	for _, key := range []string{"clusters", "contexts", "users"} {
		if items, ok := config[key].([]interface{}); ok {
			newItems, foundInSection := removeClusterFromSection(items, clusterName)
			config[key] = newItems
			if foundInSection {
				removed = true
			}
		}
	}
	return removed
}

// removeClusterFromSection removes entries with the given name from a section
func removeClusterFromSection(items []interface{}, clusterName string) ([]interface{}, bool) {
	newItems := make([]interface{}, 0, len(items))
	found := false

	for _, item := range items {
		if m, ok := item.(map[interface{}]interface{}); ok {
			if m["name"] != clusterName {
				newItems = append(newItems, item)
			} else {
				found = true
			}
		}
	}

	return newItems, found
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

func createConfigFromStdin(configPath string) error {
	// Read and validate stdin data
	stdinData, err := readAndValidateStdin()
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the stdin data to the config file
	if err := os.WriteFile(configPath, stdinData, filePerm); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Printf("Created new config file at %s from piped input\n", configPath)
	return nil
}

func addClusterConfigFromStdin(mainConfigPath string) error {
	// Read main config
	mainConfig, err := readConfig(mainConfigPath)
	if err != nil {
		return err
	}

	// Read and validate stdin data
	stdinData, err := readAndValidateStdin()
	if err != nil {
		return err
	}

	// Parse the piped config
	var newConfig map[string]interface{}
	if err := yaml.Unmarshal(stdinData, &newConfig); err != nil {
		return fmt.Errorf("invalid YAML data from stdin: %v", err)
	}

	// Create backup
	backupPath, err := backupConfig(mainConfigPath)
	if err != nil {
		return err
	}

	// Merge new config into main config with duplicate detection
	err = mergeConfigs(mainConfig, newConfig)
	if err != nil {
		return err
	}

	err = writeConfig(mainConfigPath, mainConfig)
	if err != nil {
		return err
	}

	fmt.Printf("New cluster config added from piped input. Backup created at %s\n", backupPath)
	return nil
}

// Helper function to read and validate stdin data once
func readAndValidateStdin() ([]byte, error) {
	// Read from stdin
	stdinData, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %v", err)
	}

	if len(stdinData) == 0 {
		return nil, fmt.Errorf("no data received from stdin")
	}

	// Validate that the piped data is valid YAML
	var testConfig map[string]interface{}
	if err := yaml.Unmarshal(stdinData, &testConfig); err != nil {
		return nil, fmt.Errorf("invalid YAML data from stdin: %v", err)
	}

	// Validate kubeconfig structure
	if err := validateKubeconfig(testConfig); err != nil {
		return nil, fmt.Errorf("invalid kubeconfig format: %v", err)
	}

	return stdinData, nil
}

// Validate kubeconfig structure
func validateKubeconfig(config map[string]interface{}) error {
	// Check required fields
	requiredFields := []string{"apiVersion", "kind"}
	for _, field := range requiredFields {
		if _, exists := config[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}

	// Check if it's a Config kind
	if kind, ok := config["kind"].(string); !ok || kind != "Config" {
		return fmt.Errorf("expected kind 'Config', got %v", config["kind"])
	}

	// Check if apiVersion is v1
	if apiVersion, ok := config["apiVersion"].(string); !ok || apiVersion != "v1" {
		return fmt.Errorf("expected apiVersion 'v1', got %v", config["apiVersion"])
	}

	// Ensure clusters, contexts, and users fields exist (can be empty)
	for _, field := range []string{"clusters", "contexts", "users"} {
		if _, exists := config[field]; !exists {
			config[field] = []interface{}{}
		}
	}

	return nil
}

// Improved merge function with duplicate detection
func mergeConfigs(mainConfig, newConfig map[string]interface{}) error {
	for _, key := range []string{"clusters", "contexts", "users"} {
		if newItems, ok := newConfig[key].([]interface{}); ok {
			if mainItems, ok := mainConfig[key].([]interface{}); ok {
				// Check for duplicates and warn
				duplicates := checkForDuplicates(mainItems, newItems)
				if len(duplicates) > 0 {
					fmt.Printf("Warning: Found duplicate %s names: %v. They will be added anyway.\n", key, duplicates)
				}
				mainConfig[key] = append(mainItems, newItems...)
			} else {
				mainConfig[key] = newItems
			}
		}
	}
	return nil
}

// Check for duplicate names in kubeconfig items
func checkForDuplicates(existing, newItems []interface{}) []string {
	existingNames := extractNames(existing)
	return findDuplicateNames(newItems, existingNames)
}

// extractNames extracts all names from a slice of kubeconfig items
func extractNames(items []interface{}) map[string]bool {
	names := make(map[string]bool)
	for _, item := range items {
		if name := getItemName(item); name != "" {
			names[name] = true
		}
	}
	return names
}

// findDuplicateNames finds names in newItems that already exist in existingNames
func findDuplicateNames(newItems []interface{}, existingNames map[string]bool) []string {
	var duplicates []string
	for _, item := range newItems {
		if name := getItemName(item); name != "" && existingNames[name] {
			duplicates = append(duplicates, name)
		}
	}
	return duplicates
}

// getItemName safely extracts the name from a kubeconfig item
func getItemName(item interface{}) string {
	if m, ok := item.(map[interface{}]interface{}); ok {
		if name, ok := m["name"].(string); ok {
			return name
		}
	}
	return ""
}
