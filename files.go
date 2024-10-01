package main

import (
	"bufio"
	"os"

	"github.com/rs/zerolog/log"
)

// LoadFile loads the text file `filename` and returns an array of string, the
// lines in that text file. The special filename `-` means standard input.
func LoadFile(filename string) []string {
	var lines []string = []string{}
	var fs *bufio.Scanner
	if filename != "-" {
		log.Debug().Msgf("Loading contents of file %s", filename)
		fd, err := os.Open(filename)
		if err != nil {
			log.Fatal().Msgf("Error opening file %s: %s\n", filename, err.Error())
			os.Exit(1)
		}
		defer fd.Close()
		fs = bufio.NewScanner(fd)
	} else {
		log.Debug().Msgf("Reading from standard input")
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
	var lines []string = []string{}
	for _, file := range filenames {
		lines = append(lines, LoadFile(file)...)
	}
	contents = ProcessInputFiles(lines, key)
	log.Debug().Msgf("Read %d files", len(options.files))
	return contents
}
