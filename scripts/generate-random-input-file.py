#!/usr/bin/env python3

import argparse
import random
import string


def parse_options():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--lines",
        help="The number of lines to generate",
        default=1,
        type=int,
    )
    parser.add_argument(
        "--fields",
        help="The number of fields per line",
        default=1,
        type=int,
    )
    parser.add_argument(
        "--field-length",
        help="The number of characters of a field",
        default=1,
        type=int,
    )
    return parser.parse_args()


def main():
    options = parse_options()
    for i in range(options.lines):
        line = []
        for j in range(options.fields):
            line.append("".join(random.choices(string.ascii_letters,
                                k=options.field_length)))
        print(" ".join(line))


if __name__ == "__main__":
    main()
