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
var fieldSeparator string
var files []string = []string{}
var key int
var memprofile string
var multiline int
var printVersion bool

// parseCommandLine initializes the argument parser and parses the command line.
func parseCommandLine() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTION]... [FILE]...\n", path.Clean(os.Args[0]))
		fmt.Fprintln(flag.CommandLine.Output(), `
Write sorted concatenation of all FILE(s) to standard output.

With no FILE, or when FILE is -, read standard input.

Options:`)
		flag.PrintDefaults()
	}

	flag.BoolVar(&debug, "debug", false, "Print debugging output")
	flag.BoolVar(&printVersion, "version", false, "Print version and exit")
	flag.IntVar(&key, "key", 1, "Sort lines based on a particular field")
	flag.IntVar(&multiline, "multiline", 1, "Combine multiple lines before sorting")
	flag.StringVar(&fieldSeparator, "field-separator", " ", "Use this field separator")
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

	var concatenatedContents ContentType = LoadInputFiles(files)
	var sortedContents ContentType = SortContents(concatenatedContents, key)

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	for _, multiline := range sortedContents {
		for _, line := range multiline.Lines {
			fmt.Println(line)
		}
	}
}
