package main

import (
	"bufio"
	"os"

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
