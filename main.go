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
var cpuprofile string
var debug bool
var fieldSeparator string = " "
var files []string = []string{}
var key KeyType = KeyType{}
var memprofile string
var multiline int = 1
var printVersion bool
var uniq bool

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

Options:`)
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output(), `
Key Specification:

The keys are specified with F[,[F]] where F is a field (defined by the field separator) in the multiline used to sort.

F      Use only field F for multiline comparisons
F,     Use the remainder of the multiline starting with field F for comparison
F1,F2  Use all fields between [F1,F2] for comparison
		`)
	}

	flag.BoolVar(&debug, "debug", false, "Print debugging output")
	flag.StringVar(&fieldSeparator, "field-separator", " ", "Use this field separator")
	flag.Var(&key, "key", "Sort lines based on a particular field, see 'Key Specification' for details")
	flag.IntVar(&multiline, "multiline", 1, "Combine multiple lines before sorting")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.BoolVar(&uniq, "uniq", false, "Uniq'ify the sorted multilines")

	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write cpu profile to file")
	flag.StringVar(&memprofile, "memprofile", "", "write memory profile to file")

	flag.Parse()

	if printVersion {
		fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n", Version)
		os.Exit(0)
	}
	if flag.NArg() == 0 {
		files = append(files, "-")
	} else {
		files = append(files, flag.Args()...)
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
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var concatenatedContents ContentType = LoadInputFiles(files, key)
	var sortedContents ContentType = SortContents(concatenatedContents)
	var uniqContents ContentType = UniqContents(sortedContents)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	for _, multiline := range uniqContents {
		for _, line := range multiline.Lines {
			fmt.Println(line)
		}
	}
}
