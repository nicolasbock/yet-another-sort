# Yet-another-sort

A sort util that mimics the functionality of GNU `sort`

[![Build and test](https://github.com/nicolasbock/yet-another-sort/actions/workflows/go-package.yaml/badge.svg)](https://github.com/nicolasbock/yet-another-sort/actions/workflows/go-package.yaml)

## Introduction

This project leans heavily on, and aspires to eventually reach full feature
parity with, GNU `sort`. The goal of this project is to experiment with
additional features that are not part of GNU `sort`.

The main purpose of this program is to write sorted concatenation of all
file(s) to standard output.

## Multiline Support

The user can specify how many lines `yet-another-sort` should consider as one
`multiline` unit. For example, the following `input`:

```console
$ cat input
b
1
a
2
```

can be sorted 2 lines at a time with

```console
$ yet-another-sort --multiline 2 input
a
2
b
1
```

This can be useful when sorting something like a timestamped bash history:

```console
$ cat .bash_history
#1692484702
history
#1692484723
ls -lah
#1692484726
git status
#1692484733
history
#1692484737
ls
$ cat .bash_history \
  | yet-another-sort --multiline 2 --key 2, --uniq last \
  | yet-another-sort --multiline 2 --key 1
#1692484723
ls -lah
#1692484726
git status
#1692484733
history
#1692484737
ls
```

## Installation

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/yet-another-sort)

The program can also be installed via a regular Ubuntu package:

```console
sudo add-apt-repository ppa:nicolasbock/yet-another-sort
sudo apt install yet-another-sort
```
