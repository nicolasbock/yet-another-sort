package main

import (
	"bufio"
	"os"

	"github.com/rs/zerolog/log"
)

// LoadFile loads the text file `filename` and returns an array of string, the
// lines in that text file. The special filename `-` means standard input.
func LoadFile(filename string) []string {
	var fs *bufio.Scanner
	var fd *os.File
	var err error

	if filename != "-" {
		log.Debug().Msgf("Loading contents of file %s", filename)
		fd, err = os.Open(filename)
		if err != nil {
			log.Fatal().Msgf("Error opening file %s: %s\n", filename, err.Error())
			os.Exit(1)
		}
		defer fd.Close()

		// Use larger buffer for better performance
		const bufferSize = 256 * 1024 // 256KB buffer
		reader := bufio.NewReaderSize(fd, bufferSize)
		fs = bufio.NewScanner(reader)

		// Increase scanner buffer capacity for long lines
		buf := make([]byte, bufferSize)
		fs.Buffer(buf, bufferSize)
	} else {
		log.Debug().Msgf("Reading from standard input")

		// Use larger buffer for stdin as well
		const bufferSize = 256 * 1024
		reader := bufio.NewReaderSize(os.Stdin, bufferSize)
		fs = bufio.NewScanner(reader)

		buf := make([]byte, bufferSize)
		fs.Buffer(buf, bufferSize)
	}

	fs.Split(bufio.ScanLines)

	// Pre-allocate with a reasonable initial capacity
	// This avoids multiple reallocations during append
	lines := make([]string, 0, 16384) // Start with capacity for ~16K lines

	for fs.Scan() {
		lines = append(lines, fs.Text())
	}

	if err := fs.Err(); err != nil {
		log.Fatal().Msgf("Error reading file: %s", err.Error())
	}

	return lines
}

// LoadInputFiles loads the input file(s) and returns a concatenated list of lines.
func LoadInputFiles(filenames []string, key KeyType) ContentType {
	var lines []string

	// Pre-allocate for multiple files
	if len(filenames) == 1 {
		lines = LoadFile(filenames[0])
	} else {
		// Estimate total capacity for multiple files
		lines = make([]string, 0, 16384*len(filenames))
		for _, file := range filenames {
			lines = append(lines, LoadFile(file)...)
		}
	}

	contents := ProcessInputFiles(lines, key)
	log.Debug().Msgf("Read %d files", len(filenames))
	return contents
}
