#!/bin/bash

_complete_yet_another_sort () {
  local cur prev
  local known_options=(
    --debug
    --field-separator
    --force
    --ignore-case
    --ignore-leading-blanks
    --ignore-leading-whitespace
    --key
    --multiline
    --output
    --sort-mode
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
    --sort-mode)
      readarray -t COMPREPLY < <(compgen -W 'bubble merge' -- "${cur}")
      return
      ;;
    --uniq)
      readarray -t COMPREPLY < <(compgen -W 'first last' -- "${cur}")
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
