package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// LoadFile loads the text file `filename` and returns an array of string, the
// lines in that text file. The special filename `-` means standard input.
func LoadFile(filename string) []string {
	var lines []string = []string{}
	log.Debug().Msgf("Loading contents of file %s", filename)
	var fs *bufio.Scanner
	if filename != "-" {
		fd, err := os.Open(filename)
		if err != nil {
			log.Fatal().Msgf("Error opening file %s: %s\n", filename, err.Error())
			os.Exit(1)
		}
		defer fd.Close()
		fs = bufio.NewScanner(fd)
	} else {
		fs = bufio.NewScanner(os.Stdin)
	}
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		lines = append(lines, fs.Text())
	}
	return lines
}

// ProcessInputFiles analyzes the input and splits the lines into multilines. It
// uses the key to construct the part of the multiline that will be used for
// comparisons.
func ProcessInputFiles(lines []string, key KeyType) ContentType {
	log.Debug().Msgf("Processing %d lines, key: %s", len(lines), key)
	var contents ContentType = ContentType{}
	var lastContentLine *ContentLineType
	for i := 0; i < len(lines); i += multiline {
		contents = append(contents, ContentLineType{})
		lastContentLine = &contents[len(contents)-1]
		for j := 0; j < multiline; j++ {
			var linenumber int = i + j
			if linenumber >= len(lines) {
				break
			}

			var fields []string = strings.Split(lines[linenumber], fieldSeparator)
			lastContentLine.Lines = append(lastContentLine.Lines, lines[linenumber])
			lastContentLine.Fields = append(lastContentLine.Fields, fields...)
		}

		switch key.Type {
		case NoKey:
			lastContentLine.CompareLine = strings.Join(lastContentLine.Fields, fieldSeparator)
		case SingleField:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = lastContentLine.Fields[key.Keys[0]-1]
		case Remainder:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = strings.Join(
				lastContentLine.Fields[key.Keys[0]-1:], fieldSeparator)
		case SubSet:
			if len(lastContentLine.Fields) < key.Keys[1] {
				log.Fatal().Msgf("multiline %s does not have enough keys", lastContentLine.Lines)
			}
			lastContentLine.CompareLine = strings.Join(
				lastContentLine.Fields[key.Keys[0]-1:key.Keys[1]], fieldSeparator)
		}

		if ignoreLeadingBlanks {
			lastContentLine.CompareLine = strings.TrimLeft(lastContentLine.CompareLine, " \t")
		}
	}
	return contents
}

// LoadInputFiles loads the input file(s) and returns a concatenated list of lines.
func LoadInputFiles(filenames []string, key KeyType) ContentType {
	var contents ContentType = ContentType{}
	var lines []string = []string{}
	for _, file := range filenames {
		lines = append(lines, LoadFile(file)...)
	}
	contents = ProcessInputFiles(lines, key)
	log.Debug().Msgf("Read %d files", len(files))
	return contents
}
