package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const bashCompletionScript = `
#! /bin/bash

_kubecnf_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="add remove list rollback --help -h --config -c"

    case "${prev}" in
        kubecnf)
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
        add)
            COMPREPLY=( $(compgen -f ${cur}) )
            return 0
            ;;
        remove)
            COMPREPLY=( $(compgen -W "$(kubecnf remove --generate-bash-completion)" -- ${cur}) )
            return 0
            ;;
        --config|-c)
            COMPREPLY=( $(compgen -f ${cur}) )
            return 0
            ;;
        *)
            ;;
    esac

    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
}

complete -F _kubecnf_completion kubecnf
`

func setupBashCompletion() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".config", "kubecnf")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	completionFile := filepath.Join(configDir, "kubecnf-completion.bash")
	if _, err := os.Stat(completionFile); os.IsNotExist(err) {
		err = os.WriteFile(completionFile, []byte(bashCompletionScript), 0644)
		if err != nil {
			return fmt.Errorf("failed to write completion script: %v", err)
		}

		fmt.Println("Bash completion script has been generated.")
		fmt.Printf("To enable completion, run:\n\n")
		fmt.Printf("    source %s\n\n", completionFile)
		fmt.Println("To make it permanent, add the following line to your ~/.bashrc file:")
		fmt.Printf("\n    source %s\n\n", completionFile)
	}

	return nil
}
