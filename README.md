# Yet-another-sort

A sort util that mimics the functionality of GNU `sort`

## Introction

This project leans heavily on, and aspires to eventually reach full feature
parity with GNU `sort`. The goal of this project is to experiment with
additional features that are not part of GNU `sort`.

The main purpose of this program is to write sorted concatenation of all
file(s) to standard output.

## Multiline Support

The user can specify how many lines `yet-another-sort` should consider as one
`multiline` unit. For example, using the following input:

    b
    1
    a
    2

can be sorted 2 lines at a time with

    $ yet-another-sort --multiline 2
    a
    2
    b
    1