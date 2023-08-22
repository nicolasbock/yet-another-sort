#!/bin/bash

_complete_yet_another_sort () {
  local cur prev
  local known_options=(
    --debug
    --field-separator
    --ignore-leading-blanks
    --ignore-leading-whitespace
    --key
    --multiline
    --uniq
    --version
    -h --help
  )

  _init_completion || return

  case "$prev" in
    --folder|--destination)
      _filedir
      return
      ;;
    --print-database-format)
      readarray -t COMPREPLY < <(compgen -W 'CSV JSON YAML' -- "${cur}")
      return
      ;;
    --destination-option)
      readarray -t COMPREPLY < <(compgen -W 'panic delete append' -- "${cur}")
  esac

  if [[ "$cur" == -* ]]; then
    readarray -t COMPREPLY < <(compgen -W "${known_options[*]}" -- "${cur}" )
    return
  fi

  _filedir
}

complete -F _complete_yet_another_sort yet-another-sort
