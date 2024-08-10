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
