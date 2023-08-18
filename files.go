package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// LoadFile loads the text file `filename` and returns an array of string, the
// lines in that text file. The special filename `-` means standard input.
func LoadFile(filename string) (lines []string) {
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

// LoadInputFiles loads the input file(s) and returns a concatenated list of lines.
func LoadInputFiles(filenames []string) (contents ContentType) {
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
