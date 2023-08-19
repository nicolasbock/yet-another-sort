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

// LoadInputFiles loads the input file(s) and returns a concatenated list of lines.
func LoadInputFiles(filenames []string, key KeyType) ContentType {
	var contents ContentType = ContentType{}
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
		}
		switch key.Type {
		case NoKey:
			lastContentLine.CompareLine = strings.Join(lastContentLine.Fields, fieldSeparator)
		case SingleField:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
					lineNumber, lastContentLine.Lines, file)

			}
			lastContentLine.CompareLine = lastContentLine.Fields[key.Keys[0]]
		case Remainder:
			if len(lastContentLine.Fields) < key.Keys[0] {
				log.Fatal().Msgf("multiline %d (%s) of file %s does not have enough keys",
					lineNumber, lastContentLine.Lines, file)

			}
			lastContentLine.CompareLine = strings.Join(lastContentLine.Fields[key.Keys[0]:], fieldSeparator)
		case SubSet:
			
		}
		log.Debug().Msgf("Read %d lines in file %s", lineNumber, file)
	}
	log.Debug().Msgf("Read %d files", len(files))

	return contents
}
