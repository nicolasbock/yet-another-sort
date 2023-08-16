package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var files []string = []string{}
var debug bool

var fileContents [][]string = [][]string{}

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
	flag.BoolVar(&debug, "debug", false, "Print debugging information.")
	flag.Parse()
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
func loadInputFiles(filenames []string) (contents []string) {
	for _, file := range filenames {
		log.Debug().Msgf("Loading contents of file %s", file)
		fd, err := os.Open(file)
		if err != nil {
			log.Fatal().Msgf("Error opening file %s: %s\n", file, err.Error())
			os.Exit(1)
		}
		defer fd.Close()
		fs := bufio.NewScanner(fd)
		fs.Split(bufio.ScanLines)
		for fs.Scan() {
			contents = append(contents, fs.Text())
		}
	}
	log.Debug().Msgf("Read %d files\n", len(files))

	return contents
}

// sortContents sorts the content lines and returns a sorted list.
func sortContents(contents []string) (sortedContents []string) {
	sortedContents = append(sortedContents, contents...)
	sort.Strings(sortedContents)
	return sortedContents
}

func main() {
	initializeLogging()
	parseCommandLine()
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var concatenatedContents []string = loadInputFiles(files)
	var sortedContents []string = sortContents(concatenatedContents)

	for _, line := range sortedContents {
		fmt.Printf("%s\n", line)
	}
}
