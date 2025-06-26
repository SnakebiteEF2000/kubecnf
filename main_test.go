package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "home directory expansion",
			input:    "~/test",
			expected: filepath.Join(os.Getenv("HOME"), "test"),
		},
		{
			name:     "no expansion needed",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)
			if tt.input == "~/test" {
				// For home directory expansion, just check if it doesn't start with ~
				if result[0] == '~' {
					t.Errorf("expandPath(%s) failed to expand home directory", tt.input)
				}
			} else if result != tt.expected {
				t.Errorf("expandPath(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVersionInfo(t *testing.T) {
	if version == "" {
		t.Error("version should not be empty")
	}
	if buildTime == "" {
		t.Error("buildTime should not be empty")
	}
}

func TestValidateKubeconfig(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid kubeconfig",
			config: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Config",
				"clusters":   []interface{}{},
				"contexts":   []interface{}{},
				"users":      []interface{}{},
			},
			wantErr: false,
		},
		{
			name: "missing apiVersion",
			config: map[string]interface{}{
				"kind":     "Config",
				"clusters": []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "wrong kind",
			config: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Pod",
				"clusters":   []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "wrong apiVersion",
			config: map[string]interface{}{
				"apiVersion": "v2",
				"kind":       "Config",
				"clusters":   []interface{}{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKubeconfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateKubeconfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckForDuplicates(t *testing.T) {
	existing := []interface{}{
		map[interface{}]interface{}{"name": "cluster1"},
		map[interface{}]interface{}{"name": "cluster2"},
	}

	new := []interface{}{
		map[interface{}]interface{}{"name": "cluster2"}, // duplicate
		map[interface{}]interface{}{"name": "cluster3"}, // new
	}

	duplicates := checkForDuplicates(existing, new, "clusters")

	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate, got %d", len(duplicates))
	}

	if len(duplicates) > 0 && duplicates[0] != "cluster2" {
		t.Errorf("Expected duplicate 'cluster2', got '%s'", duplicates[0])
	}
}
