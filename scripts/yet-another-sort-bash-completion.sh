#!/bin/bash

_complete_yet_another_sort () {
  local cur prev
  local known_options=(
    --cpuprofile
    --debug
    --field-separator
    --force
    --ignore-case
    --ignore-leading-blanks
    --ignore-leading-whitespace
    --key
    --memprofile
    --multiline
    --output
    --uniq
    --version
    -h --help
  )

  _init_completion || return

  case "$prev" in
    --output|--cpuprofile|--memprofile)
      _filedir
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
