#!/bin/bash

_complete_yet_another_sort () {
  local cur prev
  local known_options=(
    --debug
    --field-separator
    --force
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
    --uniq)
      readarray -t COMPREPLY < <(compgen -W 'first last' -- "${cur}")
      ;;
  esac

  if [[ "$cur" == -* ]]; then
    readarray -t COMPREPLY < <(compgen -W "${known_options[*]}" -- "${cur}" )
    return
  fi

  _filedir
}

complete -F _complete_yet_another_sort yet-another-sort
