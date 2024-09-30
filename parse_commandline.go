package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

// parseCommandLine initializes the argument parser and parses the command line.
func parseCommandLine() {
	commandLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	commandLine.Usage = func() {
		fmt.Fprintf(commandLine.Output(), "Usage: %s [OPTION]... [FILE]...\n", path.Clean(os.Args[0]))
		fmt.Fprintln(commandLine.Output(), `
Write sorted concatenation of all FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

yet-another-sort interprets the input in 'multilines' which are line groupings of one or multiple lines. This can be useful if the input contains header lines, e.g. a timestamp. Examples are the bash history file, which contains lines such as

	#1692110031
	ls
	#1692110033
	yet-another-sort --multiline 2 --key 2,

Note that when the --key option is used, multiple adjacent field-separators are treated as one separator, i.e. empty fields are ignored.

Options:`)
		fmt.Fprintln(commandLine.Output())
		commandLine.PrintDefaults()
		fmt.Fprintln(commandLine.Output(), `
Key Specification:

The keys are specified with F[,[F]] where F is a field (defined by the field separator) in the multiline used to sort.

F      Use only field F for multiline comparisons
F,     Use the remainder of the multiline starting with field F for comparison
F1,F2  Use all fields between [F1,F2] for comparison
		`)
	}

	commandLine.BoolVar(&debug, "debug", false, "Print debugging output")
	commandLine.StringVar(&fieldSeparator, "field-separator", " ", "Use this field separator")
	commandLine.BoolVar(&forceOutput, "force", false, "Overwrite output file if it exists")
	commandLine.BoolVar(&ignoreCase, "ignore-case", false, "Ignore case for comparisons")
	commandLine.BoolVar(&ignoreLeadingBlanks, "ignore-leading-blanks", false, "Ignore leading whitespace")
	commandLine.BoolVar(&ignoreLeadingBlanks, "ignore-leading-whitespace", false, "Ignore leading whitespace, same as --ignore-leading-blanks")
	commandLine.Var(&key, "key", "Sort lines based on a particular field, see 'Key Specification' for details")
	commandLine.IntVar(&multiline, "multiline", 1, "Combine multiple lines before sorting")
	commandLine.StringVar(&output, "output", "", "Write output to file instead of standard out")
	commandLine.Var(&uniq, "uniq", fmt.Sprintf("Uniq'ify the sorted multilines; keep [ \"first\", \"last\" ] of multiple identical lines; default = %s", uniq))
	commandLine.BoolVar(&printVersion, "version", false, "Print version and exit")
	commandLine.Var(&sortMode, "sort-mode", fmt.Sprintf("Choose sorting algorithm; [ \"bubble\", \"merge\" ]; default = %s", sortMode))

	commandLine.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	commandLine.StringVar(&memprofile, "memprofile", "", "write memory profile to file")

	commandLine.Parse(os.Args)

	if printVersion {
		fmt.Fprintf(commandLine.Output(), "Version: %s\n", Version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		files = append(files, "-")
	} else {
		files = append(files, commandLine.Args()...)
	}
}
