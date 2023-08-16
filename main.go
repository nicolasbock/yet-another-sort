package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Configuration options.
var debug bool
var fieldSeparator string
var files []string = []string{}
var key int
var multiline int

type ContentLineType struct {
	lines  []string
	fields []string
}

func (l ContentLineType) String() string {
	var result string
	result = "multiline\n"
	for i := range l.lines {
		result += fmt.Sprintf("  line: '%s'\n", l.lines[i])
	}
	result += fmt.Sprintf("  fields: %s\n", strings.Join(l.fields, ", "))
	return result
}

type ContentType []ContentLineType

func (c ContentType) String() string {
	var result string
	result = fmt.Sprintf("%d multilines\n", len(c))
	for _, line := range c {
		result += fmt.Sprint(line)
	}
	return result
}

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
	flag.IntVar(&key, "key", 1, "Sort lines based on a particular field")
	flag.IntVar(&multiline, "multiline", 1, "Combine multiple lines before sorting")
	flag.StringVar(&fieldSeparator, "field-separator", " ", "Use this field separator")

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
func loadInputFiles(filenames []string) (contents ContentType) {
	contents = ContentType{}
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
		var lineNumber int
		var lastContentLine *ContentLineType
		for fs.Scan() {
			if lineNumber%multiline == 0 {
				contents = append(contents, ContentLineType{})
				lastContentLine = &contents[len(contents)-1]
			}
			var line string = fs.Text()
			var fields []string = strings.Split(line, fieldSeparator)
			lastContentLine.lines = append(lastContentLine.lines, line)
			lastContentLine.fields = append(lastContentLine.fields, fields...)
			lineNumber++
			if lineNumber%multiline == 0 {
				if len(lastContentLine.fields) < key {
					log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
						lineNumber, lastContentLine.lines, file)
				}
			}
		}
		if len(contents[len(contents)-1].fields) < key {
			log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
				lineNumber, contents[len(contents)-1].lines, file)
		}
		log.Debug().Msgf("Read %d lines in file %s", lineNumber, file)
	}
	log.Debug().Msgf("Read %d files", len(files))

	return contents
}

// sortContents sorts the content lines and returns a sorted list.
func sortContents(contents ContentType) (sortedContents ContentType) {
	sortedContents = append(sortedContents, contents...)

	// Combine multilines with field-separator.
	// Sort multilined content
	// Split multilined content into original multilines

	for i := range sortedContents {
		for j := range sortedContents {
			if strings.Compare(sortedContents[i].fields[key-1], sortedContents[j].fields[key-1]) < 0 {
				var temp ContentLineType = sortedContents[i]
				sortedContents[i] = sortedContents[j]
				sortedContents[j] = temp
			}
		}
	}
	return sortedContents
}

func main() {
	initializeLogging()
	parseCommandLine()
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	var concatenatedContents ContentType = loadInputFiles(files)
	var sortedContents ContentType = sortContents(concatenatedContents)

	for _, multiline := range sortedContents {
		for _, line := range multiline.lines {
			fmt.Println(line)
		}
	}
}
