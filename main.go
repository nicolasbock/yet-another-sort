package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Program version
var Version = "unknown"

// Configuration options.
var debug bool
var fieldSeparator string
var files []string = []string{}
var key int
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

// loadInputFiles loads the input file(s) and returns a concatenated list of lines.
func loadInputFiles(filenames []string) (contents ContentType) {
	contents = ContentType{}
	for _, file := range filenames {
		var lines []string = LoadFile(file)
		var lastContentLine *ContentLineType
		var lineNumber int
		for _, line := range lines {
			if lineNumber%multiline == 0 {
				contents = append(contents, ContentLineType{})
				lastContentLine = &contents[len(contents)-1]
			}
			var fields []string = strings.Split(line, fieldSeparator)
			lastContentLine.Lines = append(lastContentLine.Lines, line)
			lastContentLine.Fields = append(lastContentLine.Fields, fields...)
			lineNumber++
			if lineNumber%multiline == 0 {
				if len(lastContentLine.Fields) < key {
					log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
						lineNumber, lastContentLine.Lines, file)
				}
			}
		}
		if len(contents[len(contents)-1].Fields) < key {
			log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
				lineNumber, contents[len(contents)-1].Lines, file)
		}
		log.Debug().Msgf("Read %d lines in file %s", lineNumber, file)
	}
	log.Debug().Msgf("Read %d files", len(files))

	return contents
}

func main() {
	initializeLogging()
	parseCommandLine()
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var concatenatedContents ContentType = loadInputFiles(files)
	var sortedContents ContentType = SortContents(concatenatedContents, key)

	for _, multiline := range sortedContents {
		for _, line := range multiline.Lines {
			fmt.Println(line)
		}
	}
}
