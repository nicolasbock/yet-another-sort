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
    --output
    --uniq
    --version
    -h --help
  )

  _init_completion || return

  case "$prev" in
    --output)
      _filedir
      return
      ;;
  esac

  if [[ "$cur" == -* ]]; then
    readarray -t COMPREPLY < <(compgen -W "${known_options[*]}" -- "${cur}" )
    return
  fi

  _filedir
}

complete -F _complete_yet_another_sort yet-another-sort
