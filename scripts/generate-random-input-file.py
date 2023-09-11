#!/usr/bin/env python3

"""Generate random input for testing."""

import argparse
import random
import string


def parse_options():
    """Parse the command line."""
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
    parser.add_argument(
        "--add-leading-whitespace",
        help="Add random numbers of leading whitespace",
        default=False,
        action="store_true"
    )
    return parser.parse_args()


def main():
    """Main function."""
    options = parse_options()
    for _ in range(options.lines):
        line = []
        for _ in range(options.fields):
            line.append("".join(random.choices(string.ascii_letters,
                                               k=options.field_length)))
        leading = ""
        if options.add_leading_whitespace:
            leading = "".join([" " * random.choice(
                [0, 0, 0, 0, 1, 2, 3])])
        print(leading + " ".join(line))


if __name__ == "__main__":
    main()
