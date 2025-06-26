package main

const bashCompletionScript = `
#! /bin/bash

_kubecnf_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="add remove list rollback completion --help -h --config -c"

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
        completion)
            COMPREPLY=( $(compgen -W "bash zsh" -- ${cur}) )
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

const zshCompletionScript = `#compdef kubecnf

_kubecnf() {
    local context curcontext="$curcontext" state line
    typeset -A opt_args

    _arguments -C \
        '(-c --config)'{-c,--config}'[specify the main config file]:config file:_files' \
        '(-h --help)'{-h,--help}'[show help]' \
        '1: :_kubecnf_commands' \
        '*:: :->args'

    case $state in
        args)
            case $words[1] in
                add)
                    _arguments \
                        '(-h --help)'{-h,--help}'[show help]' \
                        '1:input file:_files'
                    ;;
                remove)
                    _arguments \
                        '(-h --help)'{-h,--help}'[show help]' \
                        '1: :_kubecnf_clusters'
                    ;;
                completion)
                    _arguments \
                        '(-h --help)'{-h,--help}'[show help]' \
                        '1:shell:(bash zsh)'
                    ;;
                list|rollback)
                    _arguments \
                        '(-h --help)'{-h,--help}'[show help]'
                    ;;
            esac
            ;;
    esac
}

_kubecnf_commands() {
    local commands
    commands=(
        'add:add a new cluster config to the main config (from file or piped stdin)'
        'remove:remove a cluster config from the main config'
        'list:list all cluster configurations'
        'rollback:rollback to the previous config'
        'completion:output shell completion code'
        'help:Shows a list of commands or help for one command'
    )
    _describe 'command' commands
}

_kubecnf_clusters() {
    local clusters
    local config_file
    
    # Get config file from command line or use default
    config_file=${opt_args[-c]:-${opt_args[--config]:-~/.kube/config}}
    
    # Get cluster names from kubecnf, suppress errors
    clusters=(${(f)"$(kubecnf --config "$config_file" list 2>/dev/null | grep '^- ' | sed 's/^- //')"})
    
    if [[ ${#clusters[@]} -gt 0 ]]; then
        _describe 'cluster' clusters
    fi
}

_kubecnf "$@"
`
