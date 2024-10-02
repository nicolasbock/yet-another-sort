package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime/pprof"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Program version
var Version = "unknown"

// Configuration options.
type ConfigurationOptions struct {
	cpuprofile          string
	debug               bool
	fieldSeparator      string
	files               []string
	forceOutput         bool
	ignoreCase          bool
	ignoreLeadingBlanks bool
	key                 KeyType
	memprofile          string
	multiline           int
	output              string
	printVersion        bool
	sortMode            SortMode
	uniq                UniqMode
}

func NewConfigurationOptions() ConfigurationOptions {
	return ConfigurationOptions{
		fieldSeparator: " ",
		files:          []string{},
		multiline:      1,
	}
}

var options = NewConfigurationOptions()

// parseCommandLine initializes the argument parser and parses the command line.
func parseCommandLine() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTION]... [FILE]...\n", path.Clean(os.Args[0]))
		fmt.Fprintln(flag.CommandLine.Output(), `
Write sorted concatenation of all FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

yet-another-sort interprets the input in 'multilines' which are line groupings of one or multiple lines. This can be useful if the input contains header lines, e.g. a timestamp. Examples are the bash history file, which contains lines such as

	#1692110031
	ls
	#1692110033
	yet-another-sort --multiline 2 --key 2,

Note that when the --key option is used, multiple adjacent field-separators are treated as one separator, i.e. empty fields are ignored.

Options:`)
		fmt.Fprintln(flag.CommandLine.Output())
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), `
Key Specification:

The keys are specified with F[,[F]] where F is a field (defined by the field separator) in the multiline used to sort.

F      Use only field F for multiline comparisons
F,     Use the remainder of the multiline starting with field F for comparison
F1,F2  Use all fields between [F1,F2] for comparison
		`)
	}

	flag.BoolVar(&options.debug, "debug", false, "Print debugging output")
	flag.StringVar(&options.fieldSeparator, "field-separator", " ", "Use this field separator")
	flag.BoolVar(&options.forceOutput, "force", false, "Overwrite output file if it exists")
	flag.BoolVar(&options.ignoreCase, "ignore-case", false, "Ignore case for comparisons")
	flag.BoolVar(&options.ignoreLeadingBlanks, "ignore-leading-blanks", false, "Ignore leading whitespace")
	flag.BoolVar(&options.ignoreLeadingBlanks, "ignore-leading-whitespace", false, "Ignore leading whitespace, same as --ignore-leading-blanks")
	flag.Var(&options.key, "key", "Sort lines based on a particular field, see 'Key Specification' for details")
	flag.IntVar(&options.multiline, "multiline", 1, "Combine multiple lines before sorting")
	flag.StringVar(&options.output, "output", "", "Write output to file instead of standard out")
	flag.Var(&options.uniq, "uniq", fmt.Sprintf("Uniq'ify the sorted multilines; keep [ \"first\", \"last\" ] of multiple identical lines; default = %s", options.uniq))
	flag.BoolVar(&options.printVersion, "version", false, "Print version and exit")
	flag.Var(&options.sortMode, "sort-mode", fmt.Sprintf("Choose sorting algorithm; [ \"bubble\", \"merge\" ]; default = %s", options.sortMode))

	flag.StringVar(&options.cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&options.memprofile, "memprofile", "", "write memory profile to file")

	flag.Parse()

	if options.printVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", Version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		options.files = append(options.files, "-")
	} else {
		options.files = append(options.files, flag.Args()...)
	}
}

// initializeLogging initializes the logger.
func initializeLogging() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	initializeLogging()
	parseCommandLine()
	if options.debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if options.cpuprofile != "" {
		f, err := os.Create(options.cpuprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var concatenatedContents ContentType = LoadInputFiles(options.files, options.key)
	var sortedContents ContentType = SortContents(concatenatedContents)
	var uniqContents ContentType = UniqContents(sortedContents)

	if options.memprofile != "" {
		f, err := os.Create(options.memprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	var fd *os.File
	if options.output != "" {
		_, err := os.Stat(options.output)
		if err == nil && !options.forceOutput {
			log.Fatal().Msgf("output file %s already exists", options.output)
		}
		fd, err = os.Create(options.output)
		if err != nil {
			log.Fatal().Msgf("cannot open output file %s: %s", options.output, err.Error())
		}
		defer fd.Close()
	} else {
		fd = os.Stdout
	}
	for _, multiline := range uniqContents {
		for _, line := range multiline.Lines {
			fmt.Fprintln(fd, line)
		}
	}
}
